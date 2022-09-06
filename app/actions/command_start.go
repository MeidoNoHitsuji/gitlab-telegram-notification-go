package actions

import (
	"encoding/json"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab-telegram-notification-go/actions/callbacks"
	"gitlab-telegram-notification-go/telegram"
)

const Start ActionNameType = "start"

type StartAction struct {
	BaseAction
}

func NewStartAction() *StartAction {
	return &StartAction{
		BaseAction: BaseAction{
			ID:       Start,
			InitBy:   InitByCommand,
			InitText: "start",
		},
	}
}

func (act *StartAction) Active(update tgbotapi.Update) error {

	message, botMessage := telegram.GetMessageFromUpdate(update)

	if message == nil {
		return errors.New("Неизвестно откуда прилетел запрос.")
	}

	subscribeData := callbacks.NewSubscribesType()

	out, err := json.Marshal(subscribeData)

	if err != nil {
		return err
	}

	keyboard := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Подписки", string(out)),
		//TODO: Настройки
	)

	if botMessage {
		telegram.UpdateMessageById(
			message,
			"Привет! Что поделаем?",
			tgbotapi.NewInlineKeyboardMarkup(keyboard),
			nil,
		)
	} else {
		telegram.SendMessageById(
			message.Chat.ID,
			"Привет! Что поделаем?",
			tgbotapi.NewInlineKeyboardMarkup(keyboard),
			nil,
		)
	}

	return nil
}
