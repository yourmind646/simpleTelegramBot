package main

import (
	"context"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"simpleTelegramBot/fsm"
	"simpleTelegramBot/handlers"
)

func main() {
    bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
    if err != nil {
        log.Panic(err)
    }
    bot.Debug = false
	ctx := context.Background()

    log.Printf("Authorized on account %s", bot.Self.UserName)

    u := tgbotapi.NewUpdate(0)
    u.Timeout = 60

    updates := bot.GetUpdatesChan(u)
	f := fsm.NewFSM("localhost:6379", "", 0)
	handlers.LoadRoutes() // initialize all routes

	for update := range updates {
        go handlers.RouteUpdate(ctx, bot, &update, f)
    }
}
