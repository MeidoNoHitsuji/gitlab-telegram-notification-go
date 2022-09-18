package actions

import (
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab-telegram-notification-go/actions/callbacks"
)

const CreateFilter ActionNameType = "create_filter"

type CreateFilterActon struct {
	BaseAction
	CallbackData *callbacks.SelectFilterType `json:"callback_data"`
}

func NewCreateFilterActon() *CreateFilterActon {
	return &CreateFilterActon{
		BaseAction: BaseAction{
			ID:                    CreateFilter,
			InitBy:                []ActionInitByType{InitByCallback},
			InitCallbackFuncNames: []callbacks.CallbackFuncName{callbacks.CreateFilterFuncName},
			BeforeAction:          SelectFilter,
		},
	}
}

func (act *CreateFilterActon) Validate(update tgbotapi.Update) bool {
	if !act.BaseAction.Validate(update) {
		return false

	}
	return json.Unmarshal([]byte(update.CallbackQuery.Data), &act.CallbackData) == nil
}

func (act *CreateFilterActon) Active(update tgbotapi.Update) error {
	return nil
}
