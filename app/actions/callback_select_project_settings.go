package actions

import (
	"encoding/json"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/xanzy/go-gitlab"
	"gitlab-telegram-notification-go/actions/callbacks"
	"gitlab-telegram-notification-go/gitclient"
	"gitlab-telegram-notification-go/telegram"
)

const SelectProjectSettings ActionNameType = "select_project_settings"

// SelectProjectSettingsAction вызывается когда надо выбирать действие для проекта
type SelectProjectSettingsAction struct {
	BaseAction
	CallbackData *callbacks.SelectProjectSettingsType `json:"callback_data"`
	BackData     SelectProjectSettingsBackData        `json:"bd"`
}

type SelectProjectSettingsBackData struct {
	ProjectId int `json:"pi"`
}

func (act *SelectProjectSettingsAction) SetIsBack(update tgbotapi.Update) error {

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

func NewSelectProjectSettingsAction() *SelectProjectSettingsAction {
	return &SelectProjectSettingsAction{
		BaseAction: BaseAction{
			ID:                    SelectProjectSettings,
			InitBy:                []ActionInitByType{InitByCallback},
			InitCallbackFuncNames: []callbacks.CallbackFuncName{callbacks.SelectProjectSettingsFuncName},
			BeforeAction:          SelectProjectActionType,
		},
	}
}

func (act *SelectProjectSettingsAction) Validate(update tgbotapi.Update) bool {
	if !act.BaseAction.Validate(update) {
		return false

	}
	return json.Unmarshal([]byte(update.CallbackQuery.Data), &act.CallbackData) == nil
}

func (act *SelectProjectSettingsAction) Active(update tgbotapi.Update) error {
	message, _ := telegram.GetMessageFromUpdate(update)

	if message == nil {
		return errors.New("Неизвестно откуда прилетел запрос.")
	}

	var project *gitlab.Project
	var err error
	git := gitclient.Instant()

	if act.IsBack {
		project, _, err = git.Projects.GetProject(act.BackData.ProjectId, &gitlab.GetProjectOptions{})

		if err != nil {
			return err
		}
	} else {
		project, _, err = git.Projects.GetProject(act.CallbackData.ProjectId, &gitlab.GetProjectOptions{})

		if err != nil {
			return err
		}
	}

	backOut, err := json.Marshal(
		callbacks.NewBackType(
			callbacks.NewSelectProjectSettingsType(project.ID),
		),
	)

	if err != nil {
		return err
	}

	newFiletData := callbacks.NewCreateFilterType(project.ID)
	newFiletOut, err := json.Marshal(newFiletData)

	if err != nil {
		return err
	}

	selectFiletData := callbacks.NewSelectFilterType(project.ID)
	selectFiletOut, err := json.Marshal(selectFiletData)

	if err != nil {
		return err
	}

	keyboard := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Добавить", string(newFiletOut)),
		tgbotapi.NewInlineKeyboardButtonData("Обновить", string(selectFiletOut)),
	)

	keyboardBack := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Отмена", string(backOut)),
	)

	telegram.UpdateMessageById(
		message,
		"Ты выбрал настройки. Теперь выбери, хочешь ты обновить имеющийся фильтр или добавить новый?",
		tgbotapi.NewInlineKeyboardMarkup(keyboard, keyboardBack),
		nil,
	)

	return nil
}
