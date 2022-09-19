package actions

import (
	"encoding/json"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/xanzy/go-gitlab"
	"gitlab-telegram-notification-go/actions/callbacks"
	"gitlab-telegram-notification-go/database"
	"gitlab-telegram-notification-go/gitclient"
	"gitlab-telegram-notification-go/models"
	"gitlab-telegram-notification-go/telegram"
)

const EditFilter ActionNameType = "edit_filter"

type EditFilterActon struct {
	BaseAction
	CallbackData *callbacks.EditFilterType `json:"callback_data"`
	BackData     EditFilterBackData        `json:"bd"`
}

type EditFilterBackData struct {
	ProjectId int `json:"pi"`
}

func (act *EditFilterActon) SetIsBack(update tgbotapi.Update) error {

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

func NewEditFilterActon() *EditFilterActon {
	return &EditFilterActon{
		BaseAction: BaseAction{
			ID:                    EditFilter,
			InitBy:                []ActionInitByType{InitByCallback},
			InitCallbackFuncNames: []callbacks.CallbackFuncName{callbacks.EditFilterFuncName},
			BeforeAction:          SelectFilter,
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

	db := database.Instant()

	subscribeObj := models.Subscribe{
		ProjectId:         project.ID,
		TelegramChannelId: message.Chat.ID,
	}

	db.Unscoped().FirstOrCreate(&subscribeObj)

	if act.CallbackData.EventId == 0 {
		subscribeEvent := models.SubscribeEvent{
			SubscribeId: subscribeObj.ID,
			Event:       act.CallbackData.EventName,
		}

		db.Save(&subscribeEvent)
	} else {
		subscribeEvent := models.SubscribeEvent{
			ID: act.CallbackData.EventId,
		}

		result := db.Find(&subscribeEvent)

		if result.RowsAffected == 0 {
			return NewErrorForUser("Не найден передаваемый ивент для редактирования!")
		}
	}

	keyboardBack := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Отмена", string(backOut)),
	)

	telegram.UpdateMessageById(
		message,
		"А тут мы будем редактировать ивент...",
		tgbotapi.NewInlineKeyboardMarkup(keyboardBack),
		nil,
	)

	return nil
}
