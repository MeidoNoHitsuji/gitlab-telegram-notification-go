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
	fm "gitlab-telegram-notification-go/helper/formater"
	"gitlab-telegram-notification-go/models"
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
			InitCallbackFuncNames: []callbacks.CallbackFuncName{callbacks.EditFilterParameterFuncName},
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

	git := gitclient.Instant()
	_, _, err := git.Projects.GetProject(act.CallbackData.ProjectId, &gitlab.GetProjectOptions{})

	if err != nil {
		return err
	}

	db := database.Instant()

	subscribeEvent := models.SubscribeEvent{
		ID: act.CallbackData.EventId,
	}

	result := db.Find(&subscribeEvent)
	if subscribeEvent.Parameters == nil {
		subscribeEvent.Parameters = map[string][]string{}
	}

	if result.RowsAffected == 0 {
		return NewErrorForUser("Не найден передаваемый ивент для редактирования!")
	}

	var keyboards [][]tgbotapi.KeyboardButton

	allowParameters := AllowParameters()

	parameters, ok := allowParameters[subscribeEvent.Event]

	if ok {
		filters, ok := parameters[act.CallbackData.ParameterName]

		if ok {
			groupsValues := helper.Grouping(
				helper.Drop(
					helper.Keys(filters),
					AnywhereValueParameter,
				),
				3,
			)

			for _, values := range groupsValues {
				var keyboardButtons []tgbotapi.KeyboardButton
				for j := 0; j < len(values); j++ {
					keyboardButtons = append(
						keyboardButtons,
						tgbotapi.NewKeyboardButton(values[j]),
					)
				}

				keyboards = append(keyboards, tgbotapi.NewKeyboardButtonRow(keyboardButtons...))
			}
		}
	}

	keyboards = append(keyboards, tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Отмена"),
	))

	keyboard := tgbotapi.NewReplyKeyboard(keyboards...)
	keyboard.OneTimeKeyboard = true

	text := "Выберите одно из значений, предоставленных снизу или введите значение сами.\nЕсли введённое вами значение уже было фильте, то оно будет удалено."

	text += "\n\nПример ввода: " + fm.Italic("value1") + "\nИли: " + fm.Italic("value1, value2, value3")

	telegram.SendMessageById(
		message.Chat.ID,
		text,
		keyboard,
		nil,
	)

	return nil
}
