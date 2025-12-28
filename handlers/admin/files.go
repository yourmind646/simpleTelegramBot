package adminHandlers

import (
	"DeadLands/fsm"
	"DeadLands/internal/db"
	"DeadLands/internal/router"
	keyboards "DeadLands/keyboards/admin"
	"DeadLands/utils"
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mymmrac/telego"
)

func RegisterFiles(r *router.Router) {
	r.Add(
		router.Route{Check: handleAdminFilesMainChecker, Action: handleAdminFilesMain},
		router.Route{Check: handleAdminFilesBackChecker, Action: handleAdminFilesBack},
		router.Route{Check: handleAdminFilesEditChecker, Action: handleAdminFilesEdit},
		router.Route{Check: handleAdminFilesInputFileKeyChecker, Action: handleAdminFilesInputFileKey},
		router.Route{Check: handleAdminFilesSetUpChecker, Action: handleAdminFilesSetUp},
	)
}

// * files main
func handleAdminFilesMainChecker(ctx context.Context, update *telego.Update, userState string) bool {
	if update.Message == nil { // Message only
		return false
	}

	if userState == "AdminMenu.main" && update.Message.Text == keyboards.BtnFiles {
		return true
	}

	return false
}

func handleAdminFilesMain(ctx context.Context, bot *telego.Bot, update *telego.Update, f *fsm.FSM, pool *pgxpool.Pool) {
	msg_text := `<b>🗂 Файлы</b>
	
Выберите действие:`

	_, msgErr := bot.SendMessage(ctx, &telego.SendMessageParams{
		ChatID:      update.Message.Chat.ChatID(),
		Text:        msg_text,
		ParseMode:   "html",
		ReplyMarkup: keyboards.GetFilesKB(),
	})
	if msgErr != nil {
		log.Println("Ошибка отправки сообщения:", msgErr.Error())
		return
	}

	f.SetState(update.Message.From.ID, "AdminFiles", "main")
}

// * back handler
func handleAdminFilesBackChecker(ctx context.Context, update *telego.Update, userState string) bool {
	if update.Message == nil { // Message only
		return false
	}

	if strings.HasPrefix(userState, "AdminFiles.") && update.Message.Text == keyboards.BtnBack {
		return true
	}

	return false
}

func handleAdminFilesBack(ctx context.Context, bot *telego.Bot, update *telego.Update, f *fsm.FSM, pool *pgxpool.Pool) {
	handleAdminFilesMain(ctx, bot, update, f, pool)
}

// * edit files
func handleAdminFilesEditChecker(ctx context.Context, update *telego.Update, userState string) bool {
	if update.Message == nil { // Message only
		return false
	}

	if userState == "AdminFiles.main" && update.Message.Text == keyboards.BtnEdit {
		return true
	}

	return false
}

func handleAdminFilesEdit(ctx context.Context, bot *telego.Bot, update *telego.Update, f *fsm.FSM, pool *pgxpool.Pool) {
	msg_text := `Выберите ключ файла:`

	_, msgErr := bot.SendMessage(ctx, &telego.SendMessageParams{
		ChatID:      update.Message.Chat.ChatID(),
		Text:        msg_text,
		ParseMode:   "html",
		ReplyMarkup: keyboards.GetFileKeysKB(utils.AvailableFileKeys),
	})
	if msgErr != nil {
		log.Println("Ошибка отправки сообщения:", msgErr.Error())
		return
	}

	f.SetState(update.Message.From.ID, "AdminFiles", "inputFileKey")
}

// * input file key
func handleAdminFilesInputFileKeyChecker(ctx context.Context, update *telego.Update, userState string) bool {
	if update.Message == nil { // Message only
		return false
	}

	if userState == "AdminFiles.inputFileKey" && update.Message.Text != "" {
		return true
	}

	return false
}

func handleAdminFilesInputFileKey(ctx context.Context, bot *telego.Bot, update *telego.Update, f *fsm.FSM, pool *pgxpool.Pool) {
	file_key := update.Message.Text

	if !utils.IsFileKeyExists(file_key) {
		_, msgErr := bot.SendMessage(ctx, &telego.SendMessageParams{
			ChatID:      update.Message.Chat.ChatID(),
			Text:        "❌ Указанный ключ файла не найден!",
			ParseMode:   "html",
			ReplyMarkup: keyboards.GetFileKeysKB(utils.AvailableFileKeys),
		})
		if msgErr != nil {
			log.Println("Ошибка отправки сообщения:", msgErr.Error())
			return
		}
		return
	}

	qtx := db.New(pool)
	file, err := qtx.GetFileByKey(ctx, file_key)
	if err != nil {
		log.Println("Ошибка GetFileByKey:", err.Error())
		return
	}

	msg_text := `<b>Ключ:</b> %s
<b>Значение:</b> %s

Отправьте новый файл или вернитесь назад:`
	msg_text = fmt.Sprintf(msg_text, file_key, file.FileID)

	_, msgErr := bot.SendPhoto(ctx, &telego.SendPhotoParams{
		ChatID: update.Message.Chat.ChatID(),
		Photo: telego.InputFile{
			FileID: file.FileID,
		},
		Caption:     msg_text,
		ParseMode:   "html",
		ReplyMarkup: keyboards.GetBackKB(),
	})
	if msgErr != nil {
		_, msgErr = bot.SendMessage(ctx, &telego.SendMessageParams{
			ChatID:      update.Message.Chat.ChatID(),
			Text:        msg_text,
			ParseMode:   "html",
			ReplyMarkup: keyboards.GetBackKB(),
		})

		if msgErr != nil {
			log.Println("Ошибка отправки сообщения:", msgErr.Error())
			return
		}
	}

	f.UpdateData(update.Message.From.ID, map[string]string{
		"file_key": file_key,
	})
	f.SetState(update.Message.From.ID, "AdminFiles", "inputNewFileValue")
}

// * set up new file
func handleAdminFilesSetUpChecker(ctx context.Context, update *telego.Update, userState string) bool {
	if update.Message == nil { // Message only
		return false
	}

	if userState == "AdminFiles.inputNewFileValue" && len(update.Message.Photo) != 0 {
		return true
	}

	return false
}

func handleAdminFilesSetUp(ctx context.Context, bot *telego.Bot, update *telego.Update, f *fsm.FSM, pool *pgxpool.Pool) {

	data, _ := f.GetData(update.Message.From.ID)
	file_key, ok := data["file_key"]

	if !ok {
		log.Println("Не удалось получить file_key из state!")
		return
	}

	qtx := db.New(pool)
	qtx.UpdateFile(ctx, db.UpdateFileParams{
		FileID:     update.Message.Photo[len(update.Message.Photo)-1].FileID,
		FileType:   "photo",
		UploadedBy: update.Message.From.ID,
		FileKey:    file_key,
	})

	msg_text := `<b>✅ Новый файл установлен!</b>`

	_, msgErr := bot.SendMessage(ctx, &telego.SendMessageParams{
		ChatID:      update.Message.Chat.ChatID(),
		Text:        msg_text,
		ParseMode:   "html",
		ReplyMarkup: keyboards.GetFilesKB(),
	})

	if msgErr != nil {
		log.Println("Ошибка отправки сообщения:", msgErr.Error())
		return
	}

	f.SetState(update.Message.From.ID, "AdminFiles", "main")
}
