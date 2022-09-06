package actions

import (
	"encoding/json"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/xanzy/go-gitlab"
	"gitlab-telegram-notification-go/actions/callbacks"
	"gitlab-telegram-notification-go/gitclient"
	"gitlab-telegram-notification-go/telegram"
)

var SelectProjectSettings ActionNameType = "select_project_settings"

// SelectProjectSettingsAction вызывается когда надо выбирать действие для проекта
type SelectProjectSettingsAction struct {
	BaseAction
	CallbackData *callbacks.SelectProjectSettingsType
}

func NewSelectProjectSettingsAction() *SelectProjectSettingsAction {
	return &SelectProjectSettingsAction{
		BaseAction: BaseAction{
			ID:                   SelectProjectSettings,
			InitBy:               InitByCallback,
			InitCallbackFuncName: callbacks.SelectProjectSettingsFuncName,
			BeforeAction:         SelectProjectActionType,
		},
	}
}

func (act *SelectProjectSettingsAction) Active(update tgbotapi.Update) error {
	message, _ := telegram.GetMessageFromUpdate(update)

	if message == nil {
		return errors.New("Неизвестно откуда прилетел запрос.")
	}

	fmt.Println(act.CallbackData)

	git := gitclient.Instant()
	project, _, err := git.Projects.GetProject(act.CallbackData.ProjectId, &gitlab.GetProjectOptions{})

	if err != nil {
		return err
	}

	//TODO: Проверить работоспособность
	backData := callbacks.NewBackType()
	backData.BackData = SelectProjectBackData{
		projectId: project.ID,
	}

	backOut, err := json.Marshal(backData)

	keyboardBack := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Отмена", string(backOut)),
	)

	telegram.UpdateMessageById(
		message,
		"Ты выбрал настройки",
		tgbotapi.NewInlineKeyboardMarkup(keyboardBack),
		nil,
	)

	return nil
}
