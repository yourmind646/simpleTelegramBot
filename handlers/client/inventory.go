package clientHandlers

import (
	"DeadLands/fsm"
	"DeadLands/internal/db"
	"DeadLands/internal/router"
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mymmrac/telego"

	keyboards "DeadLands/keyboards/client"
)

func RegisterInventory(r *router.Router) {
	r.Add(
		router.Route{Check: handleInventoryMainChecker, Action: handleInventoryMain},
	)
}

// * start command
func handleInventoryMainChecker(ctx context.Context, update *telego.Update, userState string) bool {
	if update.Message == nil { // Message only
		return false
	}

	if update.Message.Text == "🎒 Инвентарь" {
		return true
	}

	return false
}

func handleInventoryMain(ctx context.Context, bot *telego.Bot, update *telego.Update, f *fsm.FSM, pool *pgxpool.Pool) {

	qtx := db.New(pool)
	inventoryPhotoFile, err := qtx.GetFileByKey(ctx, "inventoryPhoto")
	if err != nil {
		log.Println("Ошибка GetFileByKey:", err.Error())
		return
	}

	msg_text := "<b>🎒 Выберите категорию предметов:</b>"

	_, messageErr := bot.SendPhoto(ctx, &telego.SendPhotoParams{
		ChatID:      update.Message.Chat.ChatID(),
		Photo:       telego.InputFile{FileID: inventoryPhotoFile.FileID},
		Caption:     msg_text,
		ParseMode:   "html",
		ReplyMarkup: keyboards.GetInventoryIkb(),
	})
	if messageErr != nil {
		log.Println("Ошибка отправки сообщения '👤 Персонаж':", messageErr.Error())
	}

	f.SetState(update.Message.From.ID, "MainMenu", "main")
}
