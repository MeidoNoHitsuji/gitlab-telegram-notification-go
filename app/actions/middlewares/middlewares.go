package middlewares

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab-telegram-notification-go/database"
	"gitlab-telegram-notification-go/models"
)

type MiddleWares func(update tgbotapi.Update) bool

var (
	OnlyDM MiddleWares = OnlyDMMiddleware
)

func OnlyDMMiddleware(update tgbotapi.Update) bool {
	db := database.Instant()
	var chatId int64
	if update.CallbackQuery != nil {
		chatId = update.CallbackQuery.Message.Chat.ID
	} else if update.Message != nil {
		chatId = update.Message.Chat.ID
	} else {
		return false
	}

	var user models.User

	result := db.Where(models.User{
		TelegramChannelId: chatId,
	}).First(&user)

	return result.RowsAffected > 0
}
