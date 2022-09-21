package actions

import (
	"encoding/json"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab-telegram-notification-go/actions/callbacks"
	"gitlab-telegram-notification-go/telegram"
)

const EditFilterParameter ActionNameType = "efp_act" //edit_filter_parameter

type EditFilterParameterActon struct {
	BaseAction
	CallbackData *callbacks.EditFilterWithParameterType `json:"callback_data"`
	BackData     EditFilterParameterBackData            `json:"bd"`
}

type EditFilterParameterBackData struct {
	ProjectId int `json:"pi"`
}

func (act *EditFilterParameterActon) SetIsBack(update tgbotapi.Update) error {

	err := act.BaseAction.SetIsBack(update)

	if err != nil {
		return err
	}

	var tmp map[string]interface{}
	err = json.Unmarshal([]byte(update.CallbackQuery.Data), &tmp)

	if err != nil {
		return err
	}

	backData, ok := tmp["bd"]

	if !ok {
		return errors.New("Параметры не найдены.")
	}

	out, _ := json.Marshal(backData)
	err = json.Unmarshal(out, &act.BackData)

	if err != nil {
		return err
	}

	return nil
}

func NewEditFilterParameterActon() *EditFilterParameterActon {
	return &EditFilterParameterActon{
		BaseAction: BaseAction{
			ID:                    EditFilterParameter,
			InitBy:                []ActionInitByType{InitByCallback},
			InitCallbackFuncNames: []callbacks.CallbackFuncName{callbacks.EditFilterFuncName},
			BeforeAction:          EditFilter,
		},
	}
}

func (act *EditFilterParameterActon) Validate(update tgbotapi.Update) bool {
	if !act.BaseAction.Validate(update) {
		return false

	}
	return json.Unmarshal([]byte(update.CallbackQuery.Data), &act.CallbackData) == nil
}

func (act *EditFilterParameterActon) Active(update tgbotapi.Update) error {
	message, _ := telegram.GetMessageFromUpdate(update)

	if message == nil {
		return errors.New("Неизвестно откуда прилетел запрос.")
	}

	//git := gitclient.Instant()
	//project, _, err := git.Projects.GetProject(act.CallbackData.ProjectId, &gitlab.GetProjectOptions{})
	//
	//if err != nil {
	//	return err
	//}
	//
	//telegram.SendMessageById(
	//	message.Chat.ID,
	//	"Выберите одно из значений, предоставленных снизу или введите значение сами.",
	//	tgbotapi.NewInlineKeyboardMarkup(keyboards...),
	//	nil,
	//)

	return nil
}
