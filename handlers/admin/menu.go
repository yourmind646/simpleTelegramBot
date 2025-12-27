package adminHandlers

import (
	"context"
	"log"
	"strings"

	"DeadLands/fsm"
	clientFuncs "DeadLands/handlers/client"
	"DeadLands/internal/db"
	"DeadLands/internal/router"
	keyboards "DeadLands/keyboards/admin"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterMenu(r *router.Router) {
	r.Add(
		router.Route{Check: handleAdminMenuChecker, Action: handleAdminMenu},
		router.Route{Check: handleAdminExitChecker, Action: handleAdminExit},
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
	u := update.Message.From

	// check admin
	qtx := db.New(pool)
	_, err := qtx.IsAdminExists(ctx, u.ID)
	if err != nil {
		log.Printf("[id%d] is admin exists: %s", u.ID, err.Error())
		return
	}
	//

	msg_text := "<b>🥷🏻 Вы находитесь в админ меню</b>"

	messageConf := tgbotapi.NewMessage(update.Message.From.ID, msg_text)
	messageConf.ParseMode = "html"
	messageConf.ReplyMarkup = keyboards.GetMainKB()

	_, err = bot.Send(messageConf)
	if err != nil {
		log.Panic(err)
	}

	f.SetState(update.Message.From.ID, "AdminMenu", "main")
}

// * exit command
func handleAdminExitChecker(ctx context.Context, update *tgbotapi.Update, userState string) bool {
	if strings.HasPrefix(userState, "AdminMenu.") && update.Message.Text == keyboards.BtnExit {
		return true
	}

	return false
}

func handleAdminExit(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update, f *fsm.FSM, pool *pgxpool.Pool) {
	clientFuncs.HandleStart(ctx, bot, update, f, pool)
}
