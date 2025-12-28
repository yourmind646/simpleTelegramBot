package clientHandlers

import (
	"DeadLands/fsm"
	"DeadLands/internal/db"
	"DeadLands/internal/router"
	"DeadLands/utils"
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mymmrac/telego"

	keyboards "DeadLands/keyboards/client"

	"github.com/google/uuid"
)

func RegisterMenu(r *router.Router) {
	r.Add(
		router.Route{Check: handleStartChecker, Action: HandleStart},
	)
}

// * start command
func handleStartChecker(ctx context.Context, update *telego.Update, userState string) bool {
	if update.Message == nil { // Message only
		return false
	}

	if update.Message.Text != "/start" {
		return false
	}

	return true
}

func HandleStart(ctx context.Context, bot *telego.Bot, update *telego.Update, f *fsm.FSM, pool *pgxpool.Pool) {
	// create user & hero
	tx, beginTxErr := pool.BeginTx(ctx, pgx.TxOptions{})
	if beginTxErr != nil {
		log.Println("Ошибка создания транзакции:", beginTxErr.Error())
		return
	}

	defer func() { _ = tx.Rollback(ctx) }()

	qtx := db.New(tx)

	cuErr := qtx.CreateUser(ctx, db.CreateUserParams{
		UserID:   update.Message.From.ID,
		Username: utils.ProcessRawUsername(update.Message.From.Username),
		Fullname: utils.GetFullname(update.Message.From.FirstName, update.Message.From.LastName),
	})

	if cuErr != nil {
		log.Println("CreateUser error:", cuErr.Error())
		return
	}

	_, iheErr := qtx.IsHeroExists(ctx, update.Message.From.ID)

	if errors.Is(iheErr, pgx.ErrNoRows) {
		chErr := qtx.CreateHero(ctx, db.CreateHeroParams{
			HeroID: pgtype.UUID{Bytes: uuid.New(), Valid: true},
			UserID: update.Message.From.ID,
		})

		if chErr != nil {
			log.Println("CreateHero error:", chErr.Error())
			return
		}
	} else if iheErr != nil {
		log.Println("Ошибка проверки IsHeroExists:", iheErr.Error())
		return
	}

	commitErr := tx.Commit(ctx)
	if commitErr != nil {
		log.Println("commitErr:", commitErr)
		return
	}
	//

	msg_text := "<b>☢️ Вы находитесь в главном меню</b>"

	_, messageErr := bot.SendMessage(ctx, &telego.SendMessageParams{
		ChatID:      update.Message.Chat.ChatID(),
		Text:        msg_text,
		ParseMode:   "html",
		ReplyMarkup: keyboards.GetMainKB(),
	})
	if messageErr != nil {
		log.Println("Ошибка отправки сообщения:", messageErr.Error())
	}

	f.SetState(update.Message.From.ID, "MainMenu", "main")
}
