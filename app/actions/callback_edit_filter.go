package actions

import (
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab-telegram-notification-go/actions/callbacks"
)

const EditFilter ActionNameType = "edit_filter"

type EditFilterActon struct {
	BaseAction
	CallbackData *callbacks.EditFilterType `json:"callback_data"`
}

func NewEditFilterActon() *EditFilterActon {
	return &EditFilterActon{
		BaseAction: BaseAction{
			ID:                   EditFilter,
			InitBy:               InitByCallback,
			InitCallbackFuncName: callbacks.EditFilterFuncName,
			BeforeAction:         SelectFilter,
		},
	}
}

func (act *EditFilterActon) Validate(update tgbotapi.Update) bool {
	if !act.BaseAction.Validate(update) {
		return false

	}
	return json.Unmarshal([]byte(update.CallbackQuery.Data), &act.CallbackData) == nil
}

func (act *EditFilterActon) Active(update tgbotapi.Update) error {
	return nil
}
