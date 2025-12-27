package adminHandlers

import (
	"context"
	"fmt"
	"log"

	"DeadLands/fsm"
	"DeadLands/internal/db"
	"DeadLands/internal/router"
	keyboards "DeadLands/keyboards/admin"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterStat(r *router.Router) {
	r.Add(
		router.Route{Check: handleAdminStatChecker, Action: handleAdminStat},
	)
}

// * admin command
func handleAdminStatChecker(ctx context.Context, update *tgbotapi.Update, userState string) bool {
	if userState == "AdminMenu.main" && update.Message.Text == keyboards.BtnStat {
		return true
	}

	return false
}

func handleAdminStat(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update, f *fsm.FSM, pool *pgxpool.Pool) {
	u := update.Message.From

	// get users count
	qtx := db.New(pool)
	usersCount, err := qtx.GetUsersCount(ctx)
	if err != nil {
		log.Printf("[id%d] GetUsersCount: %s", u.ID, err.Error())
		return
	}
	//

	msg_text := `<b>📊 Статистика</b>
	
Количество пользователей: %d`
	msg_text = fmt.Sprintf(msg_text, usersCount)

	messageConf := tgbotapi.NewMessage(update.Message.From.ID, msg_text)
	messageConf.ParseMode = "html"
	messageConf.ReplyMarkup = keyboards.GetMainKB()

	_, err = bot.Send(messageConf)
	if err != nil {
		log.Println("Ошибка отправки сообщения:", err.Error())
		return
	}

	f.SetState(update.Message.From.ID, "AdminMenu", "main")
}
