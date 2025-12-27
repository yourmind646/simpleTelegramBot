package clientKeyboards

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func GetMainKB() tgbotapi.ReplyKeyboardMarkup {
	kb := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🗺 Вылазки"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("👤 Персонаж"),
			tgbotapi.NewKeyboardButton("🎒 Инвентарь"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🏚 База"),
			tgbotapi.NewKeyboardButton("🛠 Крафт"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🎁 Бонусы"),
			tgbotapi.NewKeyboardButton("🧳 Торговец"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🏆 Топ"),
			tgbotapi.NewKeyboardButton("⚙️ Настройки"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("❓ Помощь"),
		),
	)

	return kb
}
