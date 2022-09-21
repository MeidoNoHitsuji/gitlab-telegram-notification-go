package actions

import (
	"encoding/json"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/xanzy/go-gitlab"
	"gitlab-telegram-notification-go/actions/callbacks"
	"gitlab-telegram-notification-go/database"
	"gitlab-telegram-notification-go/gitclient"
	"gitlab-telegram-notification-go/helper"
	"gitlab-telegram-notification-go/models"
	"gitlab-telegram-notification-go/telegram"
)

const SelectFilter ActionNameType = "sf_act" //select_filter

type SelectFilterActon struct {
	BaseAction
	CallbackData *callbacks.SelectFilterType `json:"callback_data"`
	BackData     SelectFilterBackData        `json:"bd"`
}

type SelectFilterBackData struct {
	ProjectId int `json:"pi"`
}

func (act *SelectFilterActon) SetIsBack(update tgbotapi.Update) error {

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

	db := database.Instant()

	subscribeObj := models.Subscribe{
		ProjectId:         project.ID,
		TelegramChannelId: message.Chat.ID,
	}

	db.Unscoped().FirstOrCreate(&subscribeObj)

	var subscribeEvent []models.SubscribeEvent

	result := db.Where(models.SubscribeEvent{
		SubscribeId: subscribeObj.ID,
	}).Find(&subscribeEvent)

	var keyboardRows [][]tgbotapi.InlineKeyboardButton

	if result.RowsAffected != 0 {
		events := helper.AllowEventsWithName()

		lines := len(subscribeEvent) / 3

		if len(subscribeEvent)%3 > 0 {
			lines++
		}

		for i := 0; i < lines; i++ {
			slice := (i + 1) * 3
			if slice > len(subscribeEvent) {
				slice = len(subscribeEvent)
			}

			se := subscribeEvent[i*3 : slice]
			var keyboardButtons []tgbotapi.InlineKeyboardButton
			for j := 0; j < len(se); j++ {

				eventData := callbacks.NewEditFilterWithEventIdType(project.ID, se[j].ID)
				eventOut, err := json.Marshal(eventData)

				if err != nil {
					return err
				}

				keyboardButtons = append(
					keyboardButtons,
					tgbotapi.NewInlineKeyboardButtonData(events[se[j].Event], string(eventOut)),
				)
			}
			keyboardRows = append(keyboardRows, tgbotapi.NewInlineKeyboardRow(keyboardButtons...))
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

	keyboardRows = append(keyboardRows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Отмена", string(backOut)),
	))

	telegram.UpdateMessageById(
		message,
		"Выберите один из фильтров, чтобы его отредактировать.",
		tgbotapi.NewInlineKeyboardMarkup(keyboardRows...),
		nil,
	)

	return nil
}
