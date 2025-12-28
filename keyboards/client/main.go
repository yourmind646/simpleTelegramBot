package clientKeyboards

import (
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

func GetMainKB() *telego.ReplyKeyboardMarkup {
	kb := tu.Keyboard(
		tu.KeyboardRow(
			telego.KeyboardButton{Text: "🗺 Вылазки"},
		),
		tu.KeyboardRow(
			telego.KeyboardButton{Text: "👤 Персонаж"},
			telego.KeyboardButton{Text: "🎒 Инвентарь"},
		),
		tu.KeyboardRow(
			telego.KeyboardButton{Text: "🏚 База"},
			telego.KeyboardButton{Text: "🛠 Крафт"},
		),
		tu.KeyboardRow(
			telego.KeyboardButton{Text: "🎁 Бонусы"},
			telego.KeyboardButton{Text: "🧳 Торговец"},
		),
		tu.KeyboardRow(
			telego.KeyboardButton{Text: "🏆 Топ"},
			telego.KeyboardButton{Text: "⚙️ Настройки"},
		),
		tu.KeyboardRow(
			telego.KeyboardButton{Text: "❓ Помощь"},
		),
	).WithResizeKeyboard().WithInputFieldPlaceholder("Выберите действие...")

	return kb
}

func GetInventoryIkb() *telego.InlineKeyboardMarkup {
	kb := tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			telego.InlineKeyboardButton{
				Text:         "🍖 Еда",
				CallbackData: "inv:category:food",
			},
			telego.InlineKeyboardButton{
				Text:         "💧 Питье",
				CallbackData: "inv:category:liquid",
			},
		),
		tu.InlineKeyboardRow(
			telego.InlineKeyboardButton{
				Text:         "💊 Медицина",
				CallbackData: "inv:category:medicine",
			},
			telego.InlineKeyboardButton{
				Text:         "🧰 Материалы",
				CallbackData: "inv:category:materials",
			},
		),
		tu.InlineKeyboardRow(
			telego.InlineKeyboardButton{
				Text:         "🏹 Оружие",
				CallbackData: "inv:category:weapon",
			},
			telego.InlineKeyboardButton{
				Text:         "🛡 Броня",
				CallbackData: "inv:category:armor",
			},
		),
	)

	return kb
}
