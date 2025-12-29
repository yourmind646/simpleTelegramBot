package adminKeyboards

import (
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

const (
	BtnStat       = "📊 Статистика"
	BtnFiles      = "🗂 Файлы"
	BtnExit       = "🚪 Выйти"
	BtnEdit       = "✏️ Изменить"
	BtnBackToMenu = "↩️ Вернуться в меню"
	BtnBack       = "↩️ Назад"
)

func GetMainKB() *telego.ReplyKeyboardMarkup {
	kb := tu.Keyboard(
		tu.KeyboardRow(
			tu.KeyboardButton(BtnStat),
		),
		tu.KeyboardRow(
			tu.KeyboardButton(BtnFiles),
		),
		tu.KeyboardRow(
			tu.KeyboardButton(BtnExit),
		),
	)

	return kb.WithResizeKeyboard().WithIsPersistent()
}

func GetFilesKB() *telego.ReplyKeyboardMarkup {
	kb := tu.Keyboard(
		tu.KeyboardRow(
			tu.KeyboardButton(BtnEdit),
		),
		tu.KeyboardRow(
			tu.KeyboardButton(BtnBackToMenu),
		),
	)

	return kb.WithResizeKeyboard().WithIsPersistent()
}

func GetFileKeysKB(fileKeys []string) *telego.ReplyKeyboardMarkup {
	kb := tu.Keyboard()

	tempButtons := []telego.KeyboardButton{}
	for _, fileKey := range fileKeys { // 3 in a row
		if len(tempButtons) == 3 {
			kb.Keyboard = append(kb.Keyboard, tempButtons)
			tempButtons = []telego.KeyboardButton{
				{Text: fileKey},
			}
		} else {
			tempButtons = append(tempButtons, telego.KeyboardButton{Text: fileKey})
		}
	}
	kb.Keyboard = append(kb.Keyboard, tempButtons)

	kb.Keyboard = append(kb.Keyboard, tu.KeyboardRow(
		tu.KeyboardButton(BtnBack),
	))

	return kb.WithResizeKeyboard().WithIsPersistent()
}

func GetBackKB() *telego.ReplyKeyboardMarkup {
	kb := tu.Keyboard(
		tu.KeyboardRow(
			tu.KeyboardButton(BtnBack),
		),
	)

	return kb.WithResizeKeyboard()
}
