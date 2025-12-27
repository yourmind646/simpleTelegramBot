package admin

import (
	"context"
	"log"

	"DeadLands/fsm"
	"DeadLands/internal/db"
	"DeadLands/internal/router"
	"DeadLands/keyboards"
	"DeadLands/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterMenu(r *router.Router) {
	r.Add(
		router.Route{Check: handleAdminMenuChecker, Action: handleAdminMenu},
	)
}

// * admin command
func handleAdminMenuChecker(ctx context.Context, update *tgbotapi.Update, userState string) bool {
	if update.Message.Command() != "admin" {
		return false
	}

	return true
}

func handleAdminMenu(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update, f *fsm.FSM, pool *pgxpool.Pool) {
	// add user
	qtx := db.New(pool)
	err := qtx.CreateUser(ctx, db.CreateUserParams{
		UserID:   update.Message.From.ID,
		Username: utils.ProcessRawUsername(update.Message.From.UserName),
		Fullname: utils.GetFullname(update.Message.From.FirstName, update.Message.From.LastName),
	})
	if err != nil {
		log.Println("create user:", err)
		return
	}
	//

	msg_text := "<b>☢️ Ты находишься в главном меню</b>"

	messageConf := tgbotapi.NewMessage(update.Message.From.ID, msg_text)
	messageConf.ParseMode = "html"
	messageConf.ReplyMarkup = keyboards.GetMainKB()

	_, err = bot.Send(messageConf)
	if err != nil {
		log.Panic(err)
	}

	f.SetState(update.Message.From.ID, "MainMenu", "main")
}
