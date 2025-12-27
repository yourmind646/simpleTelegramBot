package router

import (
	"context"
	"log"
	"sync"
	"time"

	"DeadLands/fsm"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CheckFunc func(ctx context.Context, update *tgbotapi.Update, userState string) bool
type ActionFunc func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update, f *fsm.FSM, pool *pgxpool.Pool)

type Route struct {
	Check  CheckFunc
	Action ActionFunc
}

type Router struct {
	mu     sync.RWMutex
	routes []Route
}

func New() *Router {
	return &Router{}
}

func (r *Router) Add(routes ...Route) {
	r.mu.Lock()
	r.routes = append(r.routes, routes...)
	r.mu.Unlock()
}

func (r *Router) RouteUpdate(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update, f *fsm.FSM, pool *pgxpool.Pool) {
	initTime := time.Now()

	var currentState string
	currentUser := update.SentFrom()
	if currentUser != nil {
		currentState, _ = f.GetState(currentUser.ID)
	}

	r.mu.RLock()
	routes := r.routes
	r.mu.RUnlock()

	for _, route := range routes {
		if route.Check(ctx, update, currentState) {
			route.Action(ctx, bot, update, f, pool)
			break
		}
	}

	log.Printf("Апдейт обработан за %v\n", time.Since(initTime))
}
