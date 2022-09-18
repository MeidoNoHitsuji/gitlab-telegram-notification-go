package actions

import (
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab-telegram-notification-go/actions/callbacks"
)

const BackCallbackActionType ActionNameType = "back_callback"

type BackCallbackAction struct {
	BaseAction
	CallbackData *callbacks.BackType `json:"callback_data"`
}

func (act *BackCallbackAction) Validate(update tgbotapi.Update) bool {
	if !act.BaseAction.Validate(update) {
		return false

	}
	return json.Unmarshal([]byte(update.CallbackQuery.Data), &act.CallbackData) == nil
}

func NewBackCallbackAction() *BackCallbackAction {
	return &BackCallbackAction{
		BaseAction: BaseAction{
			InitBy:                []ActionInitByType{InitByCallback},
			InitCallbackFuncNames: []callbacks.CallbackFuncName{callbacks.BackFuncName},
		},
	}
}

func (act *BackCallbackAction) Active(update tgbotapi.Update) error {
	return BackAction(update)
}
