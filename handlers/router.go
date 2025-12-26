package handlers

import (
	"context"
	"log"
	"sync"
	"time"

	"simpleTelegramBot/fsm"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CheckFunc func(ctx context.Context, update *tgbotapi.Update, userState string) bool
type ActionFunc func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update, f *fsm.FSM)

type HandlerRoute struct {
	Check CheckFunc
	Action ActionFunc
}

func (hr *HandlerRoute) Apply(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update, f *fsm.FSM, userState string) bool {
	if filterResult := hr.Check(ctx, update, userState); filterResult {
		hr.Action(ctx, bot, update, f)
		return true
	} else {
		return false
	}
}

var routes []HandlerRoute
var rmtx sync.Mutex

//! TODO: load routes
func LoadRoutes() {
	LoadMainMenuFuncs(&routes, &rmtx)

	log.Printf("Загружено %d маршрутов!\n", len(routes))
}
//!

func RouteUpdate(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update, f *fsm.FSM) {
	initTime := time.Now()
	
	var currentState string

	currentUser := update.SentFrom()
	if currentUser != nil {
		currentState, _ = f.GetState(currentUser.ID)
	} else {
		currentState = ""
	}

	for _, handlerRoute := range routes {
		if res := handlerRoute.Apply(ctx, bot, update, f, currentState); res {
			break
		}
	}

	log.Printf("Апдейт обработан за %v\n", time.Since(initTime))
}