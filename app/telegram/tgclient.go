package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
)

var instant *tgbotapi.BotAPI

func New(TelegramToken string) *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(TelegramToken)

	if err != nil {
		log.Fatalf("Failed to create telegram telegram: %v", err)
	}

	bot.Debug = true

	return bot
}

func Instant() *tgbotapi.BotAPI {
	if instant == nil {
		instant = New(os.Getenv("TELEGRAM_TOKEN"))
	}

	return instant
}
