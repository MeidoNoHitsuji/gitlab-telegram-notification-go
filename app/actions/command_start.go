package actions

import (
	"encoding/json"
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

	subscribeData := callbacks.NewSubscribesType()

	out, err := json.Marshal(subscribeData)

	if err != nil {
		return err
	}

	keyboard := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Подписки", string(out)),
		//TODO: Настройки
	)

	telegram.SendMessageById(
		update.Message.Chat.ID,
		"Привет! Что поделаем?",
		tgbotapi.NewInlineKeyboardMarkup(keyboard),
		nil,
	)

	return nil
}
