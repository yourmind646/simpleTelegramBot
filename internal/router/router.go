package router

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"DeadLands/fsm"

	"github.com/jackc/pgx/v5/pgxpool"
	telego "github.com/mymmrac/telego"
)

type CheckFunc func(ctx context.Context, update *telego.Update, userState string) bool
type ActionFunc func(ctx context.Context, bot *telego.Bot, update *telego.Update, f *fsm.FSM, pool *pgxpool.Pool)

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

func (r *Router) RouteUpdate(ctx context.Context, bot *telego.Bot, update *telego.Update, f *fsm.FSM, pool *pgxpool.Pool) {
	initTime := time.Now()

	var currentState string

	if update.Message != nil {
		currentState, _ = f.GetState(update.Message.From.ID)
	} else if update.CallbackQuery != nil {
		currentState, _ = f.GetState(update.CallbackQuery.From.ID)
	} else {
		log.Println("Неподдерживаемый тип Update!")
		return
	}

	r.mu.RLock()
	routes := r.routes
	r.mu.RUnlock()

	isHandled := false
	for _, route := range routes {
		if route.Check(ctx, update, currentState) {
			route.Action(ctx, bot, update, f, pool)
			isHandled = true
			break
		}
	}

	if !isHandled {
		fmt.Printf("Апдейт не нашел обработчик! (%v)\n", time.Since(initTime))
	} else {
		log.Printf("Апдейт обработан за %v\n", time.Since(initTime))
	}
}
