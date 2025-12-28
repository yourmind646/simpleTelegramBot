package adminHandlers

import (
	"context"
	"fmt"
	"log"

	"DeadLands/fsm"
	"DeadLands/internal/db"
	"DeadLands/internal/router"
	keyboards "DeadLands/keyboards/admin"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mymmrac/telego"
)

func RegisterStat(r *router.Router) {
	r.Add(
		router.Route{Check: handleAdminStatChecker, Action: handleAdminStat},
	)
}

// * admin command
func handleAdminStatChecker(ctx context.Context, update *telego.Update, userState string) bool {
	if update.Message == nil { // Message only
		return false
	}

	if userState == "AdminMenu.main" && update.Message.Text == keyboards.BtnStat {
		return true
	}

	return false
}

func handleAdminStat(ctx context.Context, bot *telego.Bot, update *telego.Update, f *fsm.FSM, pool *pgxpool.Pool) {
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

	_, err = bot.SendMessage(ctx, &telego.SendMessageParams{
		ChatID:    update.Message.Chat.ChatID(),
		Text:      msg_text,
		ParseMode: "html",
	})
	if err != nil {
		log.Println("Ошибка отправки сообщения:", err.Error())
		return
	}

	f.SetState(update.Message.From.ID, "AdminMenu", "main")
}
