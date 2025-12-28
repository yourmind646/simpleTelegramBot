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

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mymmrac/telego"
)

func RegisterMenu(r *router.Router) {
	r.Add(
		router.Route{Check: handleAdminMenuChecker, Action: handleAdminMenu},
		router.Route{Check: handleAdminExitChecker, Action: handleAdminExit},
		router.Route{Check: handleAdminBackToMainChecker, Action: handleAdminBackToMain},
	)
}

// * admin command
func handleAdminMenuChecker(ctx context.Context, update *telego.Update, userState string) bool {
	if update.Message == nil { // Message only
		return false
	}

	if update.Message.Text != "/admin" {
		return false
	}

	return true
}

func handleAdminMenu(ctx context.Context, bot *telego.Bot, update *telego.Update, f *fsm.FSM, pool *pgxpool.Pool) {
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

	_, err = bot.SendMessage(ctx, &telego.SendMessageParams{
		ChatID:      update.Message.Chat.ChatID(),
		Text:        msg_text,
		ParseMode:   "html",
		ReplyMarkup: keyboards.GetMainKB(),
	})
	if err != nil {
		log.Panic(err)
	}

	f.SetState(update.Message.From.ID, "AdminMenu", "main")
}

// * exit command
func handleAdminExitChecker(ctx context.Context, update *telego.Update, userState string) bool {
	if update.Message == nil { // Message only
		return false
	}

	if strings.HasPrefix(userState, "AdminMenu.") && update.Message.Text == keyboards.BtnExit {
		return true
	}

	return false
}

func handleAdminExit(ctx context.Context, bot *telego.Bot, update *telego.Update, f *fsm.FSM, pool *pgxpool.Pool) {
	clientFuncs.HandleStart(ctx, bot, update, f, pool)
}

// * back to main menu
func handleAdminBackToMainChecker(ctx context.Context, update *telego.Update, userState string) bool {
	if update.Message == nil { // Message only
		return false
	}

	if strings.HasPrefix(userState, "Admin") && update.Message.Text == keyboards.BtnBackToMenu {
		return true
	}

	return false
}

func handleAdminBackToMain(ctx context.Context, bot *telego.Bot, update *telego.Update, f *fsm.FSM, pool *pgxpool.Pool) {
	handleAdminMenu(ctx, bot, update, f, pool)
}
