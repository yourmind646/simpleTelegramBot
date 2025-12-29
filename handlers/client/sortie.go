package clientHandlers

import (
	"DeadLands/fsm"
	"DeadLands/internal/db"
	"DeadLands/internal/router"
	"DeadLands/utils"
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mymmrac/telego"

	keyboards "DeadLands/keyboards/client"
)

func RegisterSortie(r *router.Router) {
	r.Add(
		router.Route{Check: handleSortieLocationsChecker, Action: handleSortieLocations},
		router.Route{Check: handleSortieLocationsBackChecker, Action: handleSortieLocationsBack},
		router.Route{Check: handleSortieChooseLocChecker, Action: handleSortieChooseLoc},
	)
}

// ! test variant
type Location struct {
	Name        string
	FileKey     string
	DurationSec int
	Drops       map[string]map[string]int // {"waterBottle": {"maxN": 3, "pct": 33}}
}

var locations map[string]Location = map[string]Location{
	"yards": Location{
		Name:        "🏘 Дворы",
		FileKey:     "locYardsPhoto",
		DurationSec: 60,
		Drops: map[string]map[string]int{
			"water_bottle": {
				"maxN": 2,
				"pct":  50,
			},
			"canned_food": {
				"maxN": 2,
				"pct":  35,
			},
			"first_aid_kit": {
				"maxN": 1,
				"pct":  8,
			},
		},
	},
	"supermarket": Location{
		Name:        "🏪 Супермаркет",
		FileKey:     "locSupermarketPhoto",
		DurationSec: 120,
		Drops: map[string]map[string]int{
			"canned_food": {
				"maxN": 4,
				"pct":  65,
			},
			"water_bottle": {
				"maxN": 3,
				"pct":  55,
			},
			"first_aid_kit": {
				"maxN": 1,
				"pct":  12,
			},
		},
	},
	"metro": Location{
		Name:        "🚇 Метро",
		FileKey:     "locMetroPhoto",
		DurationSec: 180,
		Drops: map[string]map[string]int{
			"water_bottle": {
				"maxN": 2,
				"pct":  35,
			},
			"canned_food": {
				"maxN": 2,
				"pct":  30,
			},
			"first_aid_kit": {
				"maxN": 1,
				"pct":  18,
			},
			"bow_simple": {
				"maxN": 1,
				"pct":  6,
			},
		},
	},
	"armyBase": Location{
		Name:        "🪖 Военная часть",
		FileKey:     "locArmyBasePhoto",
		DurationSec: 240,
		Drops: map[string]map[string]int{
			"first_aid_kit": {
				"maxN": 1,
				"pct":  25,
			},
			"bow_simple": {
				"maxN": 1,
				"pct":  14,
			},
			"water_bottle": {
				"maxN": 2,
				"pct":  20,
			},
			"canned_food": {
				"maxN": 2,
				"pct":  18,
			},
		},
	},
}

// * location's list
func handleSortieLocationsChecker(ctx context.Context, update *telego.Update, userState string) bool {
	if update.Message == nil { // Message only
		return false
	}

	if update.Message.Text != "🗺 Вылазки" {
		return false
	}

	return true
}

func handleSortieLocations(ctx context.Context, bot *telego.Bot, update *telego.Update, f *fsm.FSM, pool *pgxpool.Pool) {

	qtx := db.New(pool)
	sortiePhotoFile, err := qtx.GetFileByKey(ctx, "sortiePhoto")
	if err != nil {
		log.Println("Ошибка GetFileByKey:", err.Error())
		return
	}

	// подддержка обоих типов
	var chatID telego.ChatID
	if update.Message == nil {
		chatID = telego.ChatID{ID: update.CallbackQuery.From.ID}
	} else {
		chatID = telego.ChatID{ID: update.Message.From.ID}
	}

	msg_text := "<b>🗺 Выберите локацию вылазки:</b>"

	_, messageErr := bot.SendPhoto(ctx, &telego.SendPhotoParams{
		ChatID:      chatID,
		Photo:       telego.InputFile{FileID: sortiePhotoFile.FileID},
		Caption:     msg_text,
		ParseMode:   "html",
		ReplyMarkup: keyboards.GetLocationsIkb(),
	})
	if messageErr != nil {
		log.Println("Ошибка отправки сообщения:", messageErr.Error())
	}

	f.SetState(chatID.ID, "MainMenu", "main")
}

// * back to location's list
func handleSortieLocationsBackChecker(ctx context.Context, update *telego.Update, userState string) bool {
	if update.CallbackQuery == nil { // CQ only
		return false
	}

	if update.CallbackQuery.Data == "sortie:back" && userState == "MainMenu.main" {
		return true
	}

	return false
}

func handleSortieLocationsBack(ctx context.Context, bot *telego.Bot, update *telego.Update, f *fsm.FSM, pool *pgxpool.Pool) {
	err := bot.DeleteMessage(ctx, &telego.DeleteMessageParams{
		ChatID:    update.CallbackQuery.Message.GetChat().ChatID(),
		MessageID: update.CallbackQuery.Message.GetMessageID(),
	})
	if err != nil {
		log.Println("Ошибка при удалении собщения:", err.Error())
	}
	handleSortieLocations(ctx, bot, update, f, pool)
}

// * choose location
func handleSortieChooseLocChecker(ctx context.Context, update *telego.Update, userState string) bool {
	if update.CallbackQuery == nil { // Message only
		return false
	}

	if strings.HasPrefix(update.CallbackQuery.Data, "sortie:loc:") && userState == "MainMenu.main" {
		return true
	}

	return false
}

func handleSortieChooseLoc(ctx context.Context, bot *telego.Bot, update *telego.Update, f *fsm.FSM, pool *pgxpool.Pool) {

	cdata := strings.Split(update.CallbackQuery.Data, ":")
	locKey := cdata[len(cdata)-1]

	locationObj, ok := locations[locKey]
	if !ok {
		log.Println("Ошибка загрузки локации в вылазке!")
		bot.AnswerCallbackQuery(ctx, &telego.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
			ShowAlert:       true,
			Text:            "Не удалось загрузить локацию!",
		})
		return
	}
	msg_text := `<b>«%s»</b>

<b>⏳ Длительность вылазки:</b> %s

<b>🔎 Что можно найти:</b>
%s`
	msg_text = fmt.Sprintf(
		msg_text,
		locationObj.Name,
		utils.FormatSecondsToString(locationObj.DurationSec),
		utils.FormatDropList(locationObj.Drops),
	)

	qtx := db.New(pool)
	locPhotoFile, err := qtx.GetFileByKey(ctx, locationObj.FileKey)
	if err != nil {
		log.Println("Ошибка GetFileByKey:", err.Error())
		return
	}

	_, err = bot.EditMessageMedia(ctx, &telego.EditMessageMediaParams{
		ChatID:    update.CallbackQuery.Message.GetChat().ChatID(),
		MessageID: update.CallbackQuery.Message.GetMessageID(),
		Media: &telego.InputMediaPhoto{
			Type:      "photo",
			Media:     telego.InputFile{FileID: locPhotoFile.FileID},
			Caption:   msg_text,
			ParseMode: "html",
		},
		ReplyMarkup: keyboards.GetLocationPreviewIkb(),
	})
	if err != nil {
		log.Println("Ошибка изменения сообщения:", err.Error())
		return
	}

	bot.AnswerCallbackQuery(ctx, &telego.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
	})

	f.UpdateData(update.CallbackQuery.From.ID, map[string]string{
		"locKey": locKey,
	})
}
