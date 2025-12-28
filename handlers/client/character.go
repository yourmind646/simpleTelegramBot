package clientHandlers

import (
	"DeadLands/fsm"
	"DeadLands/internal/db"
	"DeadLands/internal/router"
	"context"
	"fmt"
	"html"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mymmrac/telego"

	keyboards "DeadLands/keyboards/client"
)

func RegisterCharacter(r *router.Router) {
	r.Add(
		router.Route{Check: handleCharacterMainChecker, Action: handleCharacterMain},
	)
}

// * start command
func handleCharacterMainChecker(ctx context.Context, update *telego.Update, userState string) bool {
	if update.Message == nil { // Message only
		return false
	}

	if update.Message.Text == "👤 Персонаж" {
		return true
	}

	return false
}

func handleCharacterMain(ctx context.Context, bot *telego.Bot, update *telego.Update, f *fsm.FSM, pool *pgxpool.Pool) {

	qtx := db.New(pool)
	hero, err := qtx.GetHeroByUser(ctx, update.Message.From.ID)
	if err != nil {
		log.Println("Ошибка GetHeroByUser:", err.Error())
		return
	}
	user, err := qtx.GetUser(ctx, update.Message.From.ID)
	if err != nil {
		log.Println("Ошибка GetUser:", err.Error())
		return
	}
	heroPhotoFile, err := qtx.GetFileByKey(ctx, "profilePhoto")
	if err != nil {
		log.Println("Ошибка GetFileByKey:", err.Error())
		return
	}

	msg_text := fmt.Sprintf(
		`<b>👤 Персонаж «%s»</b>

Состояние:
❤️ Здоровье: <b>%d</b>/100
⚡️ Энергия: <b>%d</b>/100
🍖 Голод: <b>%d</b>/100
💧 Жажда: <b>%d</b>/100
☢️ Радиация: <b>%d</b>/100`,
		html.EscapeString(user.Fullname.String),
		hero.Hp, hero.Energy, hero.Hunger, hero.Thirst, hero.Radiation,
	)

	_, messageErr := bot.SendPhoto(ctx, &telego.SendPhotoParams{
		ChatID:      update.Message.Chat.ChatID(),
		Photo:       telego.InputFile{FileID: heroPhotoFile.FileID},
		Caption:     msg_text,
		ParseMode:   "html",
		ReplyMarkup: keyboards.GetMainKB(),
	})
	if messageErr != nil {
		log.Println("Ошибка отправки сообщения '👤 Персонаж':", messageErr.Error())
	}

	f.SetState(update.Message.From.ID, "MainMenu", "main")
}
