package actions

import (
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab-telegram-notification-go/actions/callbacks"
)

const SelectFilter ActionNameType = "select_filter"

type SelectFilterActon struct {
	BaseAction
	CallbackData *callbacks.SelectFilterType `json:"callback_data"`
}

func NewSelectFilterActon() *SelectFilterActon {
	return &SelectFilterActon{
		BaseAction: BaseAction{
			ID:                    SelectFilter,
			InitBy:                []ActionInitByType{InitByCallback},
			InitCallbackFuncNames: []callbacks.CallbackFuncName{callbacks.SelectFilterFuncName},
			BeforeAction:          SelectProjectSettings,
		},
	}
}

func (act *SelectFilterActon) Validate(update tgbotapi.Update) bool {
	if !act.BaseAction.Validate(update) {
		return false

	}
	return json.Unmarshal([]byte(update.CallbackQuery.Data), &act.CallbackData) == nil
}

func (act *SelectFilterActon) Active(update tgbotapi.Update) error {
	return nil
}
