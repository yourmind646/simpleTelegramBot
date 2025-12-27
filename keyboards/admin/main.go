package adminKeyboards

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

const (
	BtnStat  = "📊 Статистика"
	BtnFiles = "🗂 Файлы"
	BtnExit  = "🚪 Выйти"
)

func GetMainKB() tgbotapi.ReplyKeyboardMarkup {
	kb := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(BtnStat),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(BtnFiles),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(BtnExit),
		),
	)

	return kb
}
