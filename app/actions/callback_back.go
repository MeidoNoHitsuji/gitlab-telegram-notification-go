package actions

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab-telegram-notification-go/actions/callbacks"
)

const BackCallbackActionType ActionNameType = "back_callback"

type BackCallbackAction struct {
	BaseAction
	CallbackData *callbacks.BackType
}

func NewBackCallbackAction() *BackCallbackAction {
	return &BackCallbackAction{
		BaseAction: BaseAction{
			InitBy:               InitByCallback,
			InitCallbackFuncName: callbacks.BackFuncName,
		},
	}
}

func (act *BackCallbackAction) Active(update tgbotapi.Update) error {
	return BackAction(update)
}
