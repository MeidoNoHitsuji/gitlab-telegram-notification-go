package telegram

import (
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab-telegram-notification-go/database"
	"gitlab-telegram-notification-go/models"
	"log"
)

func SendMessage(channel *models.TelegramChannel, message string, keyboard interface{}, entities []tgbotapi.MessageEntity) (*tgbotapi.Message, error) {
	bot := Instant()

	if !channel.Active {
		err := fmt.Sprintf("Чат с Id %d недоступен!", channel.ID)
		log.Println(err)
		return nil, errors.New(err)
	}

	msgConf := tgbotapi.NewMessage(channel.ID, message)
	msgConf.ParseMode = tgbotapi.ModeMarkdown
	msgConf.DisableWebPagePreview = true
	msgConf.Entities = entities
	
	if keyboard != nil {
		msgConf.ReplyMarkup = keyboard
	}

	if channel.ID < 0 {
		msgConf.DisableNotification = true
	}

	msg, err := bot.Send(msgConf)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &msg, nil
}

func SendMessageById(telegramId int64, message string, keyboard interface{}, entities []tgbotapi.MessageEntity) (*tgbotapi.Message, error) {
	db := database.Instant()

	var channel models.TelegramChannel

	result := db.Where(&models.TelegramChannel{ID: telegramId}).Find(&channel)

	if result.Error != nil {
		log.Println(result.Error)
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		err := fmt.Sprintf("Канал с Id равному %d не найден.", telegramId)
		log.Println(err)

		return nil, errors.New(err)
	}

	return SendMessage(&channel, message, keyboard, entities)
}

func SendMessageByUsername(username string, message string, keyboard interface{}, entities []tgbotapi.MessageEntity) (*tgbotapi.Message, error) {
	db := database.Instant()

	var user models.User
	result := db.Where(&models.User{Username: username}).Preload("TelegramChannel").Find(&user)

	if result.Error != nil {
		log.Println(result.Error)
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		err := fmt.Sprintf("Пользователь с Username равному %s не найден.", username)
		log.Println(err)

		return nil, errors.New(err)
	}

	return SendMessage(&user.TelegramChannel, message, keyboard, entities)
}
