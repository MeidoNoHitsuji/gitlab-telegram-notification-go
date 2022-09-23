package actions

import (
	"encoding/json"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/xanzy/go-gitlab"
	"gitlab-telegram-notification-go/actions/callbacks"
	"gitlab-telegram-notification-go/gitclient"
	"gitlab-telegram-notification-go/helper"
	"gitlab-telegram-notification-go/telegram"
)

const CreateFilter ActionNameType = "cf_act" //create_filter

type CreateFilterActon struct {
	BaseAction
	CallbackData *callbacks.SelectFilterType `json:"callback_data"`
	BackData     CreateFilterBackData        `json:"bd"`
}

type CreateFilterBackData struct {
	ProjectId int `json:"pi"`
}

func (act *CreateFilterActon) SetIsBack(update tgbotapi.Update) error {

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

func NewCreateFilterActon() *CreateFilterActon {
	return &CreateFilterActon{
		BaseAction: BaseAction{
			ID:                    CreateFilter,
			InitBy:                []ActionInitByType{InitByCallback},
			InitCallbackFuncNames: []callbacks.CallbackFuncName{callbacks.CreateFilterFuncName},
			BeforeAction:          SelectProjectSettings,
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

	events := helper.AllowEventsWithName()
	eventKeys := helper.Keys(events)

	var keyboardRows [][]tgbotapi.InlineKeyboardButton

	for _, keys := range helper.Grouping(eventKeys, 3) {
		var keyboardButtons []tgbotapi.InlineKeyboardButton
		for j := 0; j < len(keys); j++ {

			eventData := callbacks.NewEditFilterWithEventNameType(project.ID, keys[j])
			eventOut, err := json.Marshal(eventData)

			if err != nil {
				return err
			}

			keyboardButtons = append(
				keyboardButtons,
				tgbotapi.NewInlineKeyboardButtonData(events[keys[j]], string(eventOut)),
			)
		}

		keyboardRows = append(keyboardRows, tgbotapi.NewInlineKeyboardRow(keyboardButtons...))
	}

	backOut, err := json.Marshal(
		callbacks.NewBackType(
			callbacks.NewSelectProjectSettingsType(project.ID),
		),
	)

	if err != nil {
		return err
	}

	keyboardRows = append(keyboardRows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Отмена", string(backOut)),
	))

	telegram.UpdateMessageById(
		message,
		"И так, давай создадим новый фильтр. Для этого тебе нужно выбрать ивент, на который данный фильтр будет триггерится.",
		tgbotapi.NewInlineKeyboardMarkup(keyboardRows...),
		nil,
	)

	return nil
}
