package client

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
)

var telegramInstant *tgbotapi.BotAPI

func Telegram() *tgbotapi.BotAPI {
	if telegramInstant != nil {
		return telegramInstant
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_TOKEN"))

	if err != nil {
		log.Fatalf("Failed to create telegram client: %v", err)
	}

	bot.Debug = true

	telegramInstant = bot
	return telegramInstant
}
