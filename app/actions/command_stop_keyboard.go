package actions

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab-telegram-notification-go/telegram"
)

const StopKeyboard ActionNameType = "sk"

type StopKeyboardAction struct {
	BaseAction
}

func NewStopKeyboardAction() *StopKeyboardAction {
	return &StopKeyboardAction{
		BaseAction: BaseAction{
			ID:       StopKeyboard,
			InitBy:   []ActionInitByType{InitByCommand},
			InitText: "stop_keyboard",
		},
	}
}

func (act *StopKeyboardAction) Active(update tgbotapi.Update) error {
	message, _ := telegram.GetMessageFromUpdate(update)
	telegram.SendRemoveKeyboard(message.Chat.ID, false)
	return nil
}
