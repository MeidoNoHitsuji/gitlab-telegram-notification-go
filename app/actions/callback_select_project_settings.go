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

	git := gitclient.Instant()
	project, _, err := git.Projects.GetProject(act.CallbackData.ProjectId, &gitlab.GetProjectOptions{})

	if err != nil {
		return err
	}

	backOut, err := json.Marshal(
		callbacks.NewBackType(
			callbacks.NewSelectProjectSettingsType(project.ID),
		),
	)

	if err != nil {
		return err
	}

	keyboardBack := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Отмена", string(backOut)),
	)

	telegram.UpdateMessageById(
		message,
		"Ты выбрал настройки. Теперь выбери, хочешь ты обновить имеющийся фильтр или добавить новый?",
		tgbotapi.NewInlineKeyboardMarkup(keyboardBack),
		nil,
	)

	return nil
}
