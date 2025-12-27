package main

import (
	"context"
	"errors"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5"
	pgxpool "github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"DeadLands/fsm"
	adminHandlers "DeadLands/handlers/admin"
	clientHandlers "DeadLands/handlers/client"
	"DeadLands/internal/db"
	"DeadLands/internal/router"
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

		rawCommand = strings.Replace(rawCommand, "\n    ", "", -1)
		rawCommand = strings.Replace(rawCommand, "\n", "", -1)
		commands = append(commands, strings.TrimSpace(rawCommand))
	}

	tx, err := pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		log.Panic(err)
	}
	defer tx.Rollback(ctx)

	for _, commandSql := range commands {
		_, err = tx.Exec(ctx, commandSql)
		if err != nil {
			log.Panic(err)
		}
	}

	// добавить базового админа
	qtx := db.New(tx)
	base_admin_id, err := strconv.ParseInt(os.Getenv("BASE_ADMIN"), 10, 64)

	if err == nil {
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
		log.Println("Ошибка парсинга BaseAdmin из env!")
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

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = false
	ctx := context.Background()

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
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

	clientHandlers.RegisterMenu(r)

	for update := range updates {
		go r.RouteUpdate(ctx, bot, &update, f, pool)
	}
}
