package keyboards

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func GetMainKB() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.KeyboardButton{Text: "⚒️ Go Farm"},
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.KeyboardButton{Text: "🎒 Inventory"},
			tgbotapi.KeyboardButton{Text: "👤 Profile"},
		),
	)
}