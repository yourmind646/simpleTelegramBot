package handlers

import (
	"context"
	"log"
	"simpleTelegramBot/fsm"
	"simpleTelegramBot/keyboards"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

//* loader
func LoadMainMenuFuncs(routes *[]HandlerRoute, rmtx *sync.Mutex) {
	rmtx.Lock() // lock routes
	newRoutes := []HandlerRoute{
		{handleStartChecker, handleStart},
		{handleHelpCommandChecker, handleHelpCommand},
	}

	*routes = append(*routes, newRoutes...)
	rmtx.Unlock() // unlock routes
}

//* start command
func handleStartChecker(ctx context.Context, update *tgbotapi.Update, userState string) bool {	
	if update.Message.Command() != "start" {
		return false
	}
	
	return true
}

func handleStart(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update, f *fsm.FSM) {
	msg_text := "<b>🏘 You are in the main menu.</b>"

	messageConf := tgbotapi.NewMessage(update.Message.From.ID, msg_text)
	messageConf.ParseMode = "html"
	messageConf.ReplyMarkup = keyboards.GetMainKB()

	_, err := bot.Send(messageConf)
	if err != nil {
		log.Panic(err)
	}

	f.SetState(update.Message.From.ID, "mainMenu")
}

//* help command
func handleHelpCommandChecker(ctx context.Context, update *tgbotapi.Update, userState string) bool {	
	if update.Message.Command() != "help" {
		return false
	}
	
	return true
}

func handleHelpCommand(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update, f *fsm.FSM) {
	msg_text := "<b>❓ This is a test bot written in Go.</b>"

	messageConf := tgbotapi.NewMessage(update.Message.From.ID, msg_text)
	messageConf.ParseMode = "html"

	_, err := bot.Send(messageConf)
	if err != nil {
		log.Panic(err)
	}
}