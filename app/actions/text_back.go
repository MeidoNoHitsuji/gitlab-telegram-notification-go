package actions

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const BackTextActionType ActionNameType = "back_text"

type BackTextAction struct {
	BaseAction
}

func NewBackTextAction() *BackTextAction {
	return &BackTextAction{
		BaseAction: BaseAction{
			InitBy:   InitByText,
			InitText: "Отмена",
		},
	}
}

func (act *BackTextAction) Active(update tgbotapi.Update) error {
	return BackAction(update)
}
