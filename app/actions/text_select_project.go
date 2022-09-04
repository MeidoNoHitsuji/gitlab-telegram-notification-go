package actions

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab-telegram-notification-go/telegram"
)

const SelectProjectActionType ActionNameType = "select_project"

type SelectProjectAction struct {
	BaseAction
}

func (act *SelectProjectAction) Active(update tgbotapi.Update) error {
	var message *tgbotapi.Message

	if update.Message != nil {
		message = update.Message
	} else if update.CallbackQuery != nil {
		message = update.CallbackQuery.Message
	} else {
		return errors.New("Неизвестно откуда прилетел запрос.")
	}

	telegram.UpdateMessageById(
		message,
		"kekw",
		tgbotapi.NewRemoveKeyboard(false),
		nil,
	)
	return nil
}

func NewSelectProjectAction() *SelectProjectAction {
	return &SelectProjectAction{
		BaseAction: BaseAction{
			ID:          SelectProjectActionType,
			InitBy:      InitByText,
			AfterAction: SubscribesActionType,
		},
	}
}
