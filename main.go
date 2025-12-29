package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	pgxpool "github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	telego "github.com/mymmrac/telego"

	"DeadLands/fsm"
	adminHandlers "DeadLands/handlers/admin"
	clientHandlers "DeadLands/handlers/client"
	"DeadLands/internal/db"
	"DeadLands/internal/router"
	"DeadLands/utils"
)

func initDatabaseSchema(ctx context.Context, pool *pgxpool.Pool) {
	schemaBytes, err := os.ReadFile("./sql/schema.sql")
	if err != nil {
		log.Panic(err)
	}
	schemaStr := string(schemaBytes)
	commandsRaw := strings.Split(schemaStr, ";")
	commands := []string{}
	for _, rawCommand := range commandsRaw {
		if rawCommand == "" {
			continue
		}

		rawCommand = strings.TrimSpace(rawCommand)
		if rawCommand == "" {
			continue
		}
		commands = append(commands, rawCommand)
	}

	for _, commandSql := range commands {
		_, err = pool.Exec(ctx, commandSql)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				if pgErr.Code == "42710" && strings.HasPrefix(strings.ToUpper(strings.TrimSpace(commandSql)), "CREATE TYPE ") {
					continue
				}
			}
			log.Panicf("%v (sql: %q)", err, commandSql)
		}
	}

	tx, err := pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		log.Panic(err)
	}
	defer tx.Rollback(ctx)

	// добавить базового админа
	qtx := db.New(tx)
	base_admin_id, envErr := strconv.ParseInt(os.Getenv("BASE_ADMIN"), 10, 64)

	if envErr != nil {
		log.Println("Ошибка парсинга BaseAdmin из env!")
		return
	}

	uid, iueErr := qtx.IsUserExists(ctx, base_admin_id)
	if errors.Is(iueErr, pgx.ErrNoRows) {
		cuErr := qtx.CreateUser(ctx, db.CreateUserParams{
			UserID:   base_admin_id,
			Username: pgtype.Text{String: "admin", Valid: true},
			Fullname: pgtype.Text{String: "admin", Valid: true},
		})

		if cuErr != nil {
			log.Println("CreateUser error:", cuErr.Error())
			return
		}
	} else if uid > 0 {
		_, err = qtx.IsAdminExists(ctx, base_admin_id)
		if errors.Is(err, pgx.ErrNoRows) {
			err = qtx.CreateAdmin(ctx, base_admin_id)

			if err != nil {
				log.Println("Ошибка создания базового админа:", err.Error())
			}
		} else if err == nil {
			log.Println("Базовый админ уже создан.")
		} else {
			log.Println("Не удалось создать базового админа:", err.Error())
		}
	} else {
		log.Println("Ошибка проверки IsUserExists:", iueErr.Error())
		return
	}

	// инициализация файлов
	for _, fileKey := range utils.AvailableFileKeys {
		_, err = qtx.IsFileExistsByKey(ctx, fileKey)
		if errors.Is(err, pgx.ErrNoRows) {
			err = qtx.CreateFile(ctx, db.CreateFileParams{
				FileID:     "",
				FileKey:    fileKey,
				FileType:   "undefined",
				UploadedBy: base_admin_id,
			})
			if err != nil {
				log.Panic("Ошибка CreateFile:", err.Error())
			}
		} else if err != nil {
			log.Panic("Ошибка IsFileExistsByKey:", err.Error())
		}
	}

	if err := tx.Commit(ctx); err != nil {
		log.Panic(err)
	}

	log.Println("Схема БД инициализирована!")
}

func main() {
	// грузит .env из текущей рабочей директории
	if err := godotenv.Load(); err != nil {
		log.Println(".env not loaded:", err)
	}

	bot, err := telego.NewBot(os.Getenv("TELEGRAM_APITOKEN"), telego.WithDiscardLogger())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	ctx := context.Background()

	// Call method getMe (https://core.telegram.org/bots/api#getme)
	botUser, err := bot.GetMe(context.Background())
	if err != nil {
		fmt.Println("Error:", err)
	}

	// Print Bot information
	log.Printf("Инициализирован бот: @%s (id%d)\n", botUser.Username, botUser.ID)

	updates, _ := bot.UpdatesViaLongPolling(context.Background(), nil)
	f := fsm.NewFSM("localhost:6379", "", 0)

	cfg, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	// pg conf
	cfg.MaxConns = 20
	cfg.MinConns = 2
	cfg.MaxConnLifetime = time.Hour
	cfg.MaxConnIdleTime = 30 * time.Minute
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	// check connection
	if err := pool.Ping(ctx); err != nil {
		panic(err)
	}

	initDatabaseSchema(ctx, pool) // initialize db schema

	r := router.New()
	adminHandlers.RegisterMenu(r)
	adminHandlers.RegisterStat(r)
	adminHandlers.RegisterFiles(r)

	clientHandlers.RegisterMenu(r)
	clientHandlers.RegisterCharacter(r)
	clientHandlers.RegisterInventory(r)
	clientHandlers.RegisterSortie(r)

	for update := range updates {
		go r.RouteUpdate(ctx, bot, &update, f, pool)
	}
}
