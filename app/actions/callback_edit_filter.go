package actions

import (
	"encoding/json"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/xanzy/go-gitlab"
	"gitlab-telegram-notification-go/actions/callbacks"
	"gitlab-telegram-notification-go/database"
	"gitlab-telegram-notification-go/gitclient"
	"gitlab-telegram-notification-go/helper"
	"gitlab-telegram-notification-go/models"
	"gitlab-telegram-notification-go/telegram"
	"strings"
)

const EditFilter ActionNameType = "ef_act" //edit_filter

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
			ID: EditFilter,
			InitBy: []ActionInitByType{
				InitByCallback,
				InitByText,
			},
			InitCallbackFuncNames: []callbacks.CallbackFuncName{callbacks.EditFilterFuncName},
			BeforeAction:          SelectProjectSettings,
		},
	}
}

func (act *EditFilterActon) Validate(update tgbotapi.Update) bool {
	if !act.BaseAction.Validate(update) {
		return false

	}

	if act.InitializedBy == InitByText {
		message, _ := telegram.GetMessageFromUpdate(update)

		action := database.GetUserActionInChannel(
			message.Chat.ID,
			message.From.UserName,
		)

		result := json.Unmarshal([]byte(action.Parameters), &act.CallbackData) == nil

		action.Parameters = ""

		db := database.Instant()
		db.Save(action)

		return result
	} else {
		return json.Unmarshal([]byte(update.CallbackQuery.Data), &act.CallbackData) == nil
	}
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

	db := database.Instant()

	subscribeObj := models.Subscribe{
		ProjectId:         project.ID,
		TelegramChannelId: message.Chat.ID,
	}

	db.Unscoped().FirstOrCreate(&subscribeObj)

	var subscribeEvent models.SubscribeEvent

	if act.CallbackData.EventId == 0 {
		subscribeEvent = models.SubscribeEvent{
			SubscribeId: subscribeObj.ID,
			Event:       act.CallbackData.EventName,
			Parameters:  map[string][]string{},
		}

		db.Save(&subscribeEvent)
	} else {
		subscribeEvent = models.SubscribeEvent{
			ID: act.CallbackData.EventId,
		}

		result := db.Find(&subscribeEvent)

		if result.RowsAffected == 0 {
			return NewErrorForUser("Не найден передаваемый ивент для редактирования!")
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

	keyboardBack := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Вернуться к настройкам", string(backOut)),
	)

	allowParameters := AllowParameters()

	allowEvents := helper.AllowEventsWithName()
	parameters, ok := allowParameters[subscribeEvent.Event]
	var parameterKeys []string

	for k := range parameters {
		parameterKeys = append(parameterKeys, k)
	}

	var text string
	var keyboards [][]tgbotapi.InlineKeyboardButton

	if ok {
		parameterNames := ParameterNames()
		text = fmt.Sprintf("Тип ивента: %s\n—————\nИмеющиеся параметры:", allowEvents[subscribeEvent.Event])

		if act.CallbackData.ParameterName != "" {
			if act.CallbackData.DeleteValue {
				subscribeEvent.Parameters[act.CallbackData.ParameterName] = []string{}
				db.Save(&subscribeEvent)
			} else if act.CallbackData.ParameterValue != "" {
				ps := subscribeEvent.Parameters[act.CallbackData.ParameterName]
				if helper.Contains(ps, act.CallbackData.ParameterValue) {
					ps = helper.Drop(ps, act.CallbackData.ParameterValue)
				} else {
					ps = append(ps, act.CallbackData.ParameterValue)
				}
				subscribeEvent.Parameters[act.CallbackData.ParameterName] = ps
				db.Save(&subscribeEvent)
			}
		}

		if len(subscribeEvent.Parameters) == 0 {
			text = fmt.Sprintf("%s\nОтсутствуют.", text)
		} else {
			for _, key := range parameterKeys {
				currentParams, ok := subscribeEvent.Parameters[key]

				if ok && len(currentParams) > 0 {
					text = fmt.Sprintf("%s\n> %s: %s", text, parameterNames[key], strings.Join(currentParams, ", "))
				}
			}
		}

		if act.CallbackData.ParameterName != "" {
			text = fmt.Sprintf("%s\n\nРедактируемый параметр: %s", text, parameterNames[act.CallbackData.ParameterName])
			filters := parameters[act.CallbackData.ParameterName]
			filterKeys := helper.Keys(filters)
			if helper.Contains(filterKeys, AnywhereValueParameter) {

				//TODO: Добавить переход на кнопочный экшн с вводом параметров
			} else {
				for _, keys := range helper.Grouping(filterKeys, 3) {
					var keyboardButtons []tgbotapi.InlineKeyboardButton
					for j := 0; j < len(keys); j++ {

						eventData := callbacks.NewEditFilterWithParameterValueType(project.ID, subscribeEvent.ID, act.CallbackData.ParameterName, keys[j])
						eventOut, err := json.Marshal(eventData)

						if err != nil {
							return err
						}

						keyboardButtons = append(
							keyboardButtons,
							tgbotapi.NewInlineKeyboardButtonData(filters[keys[j]], string(eventOut)),
						)
					}

					keyboards = append(keyboards, tgbotapi.NewInlineKeyboardRow(keyboardButtons...))
				}
			}

			deleteParamData := callbacks.NewEditFilterWithDeleteParameterType(project.ID, subscribeEvent.ID, act.CallbackData.ParameterName)
			deleteParamOut, err := json.Marshal(deleteParamData)

			if err != nil {
				return err
			}

			dropParamData := callbacks.NewEditFilterWithEventIdType(project.ID, subscribeEvent.ID)
			dropParamOut, err := json.Marshal(dropParamData)

			if err != nil {
				return err
			}

			keyboards = append(keyboards,
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("Отчистить параметер", string(deleteParamOut)),
					tgbotapi.NewInlineKeyboardButtonData("Изменить параметр", string(dropParamOut)),
				),
			)
		} else {
			for _, keys := range helper.Grouping(parameterKeys, 3) {
				var keyboardButtons []tgbotapi.InlineKeyboardButton
				for j := 0; j < len(keys); j++ {

					eventData := callbacks.NewEditFilterWithParameterType(project.ID, subscribeEvent.ID, keys[j])
					eventOut, err := json.Marshal(eventData)

					if err != nil {
						return err
					}

					keyboardButtons = append(
						keyboardButtons,
						tgbotapi.NewInlineKeyboardButtonData(parameterNames[keys[j]], string(eventOut)),
					)
				}

				keyboards = append(keyboards, tgbotapi.NewInlineKeyboardRow(keyboardButtons...))
			}

			deleteOut, err := json.Marshal(
				callbacks.NewBackType(
					callbacks.NewSelectProjectSettingsWithDeleteEventType(project.ID, subscribeEvent.ID),
				),
			)

			if err != nil {
				return err
			}

			keyboardDelete := tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Удалить ивент", string(deleteOut)),
			)

			keyboards = append(keyboards, keyboardDelete)
		}
	} else {
		text = "Данный ивент ещё нельзя редактировать!"
	}

	keyboards = append(keyboards, keyboardBack)

	telegram.UpdateMessageById(
		message,
		text,
		tgbotapi.NewInlineKeyboardMarkup(keyboards...),
		nil,
	)

	return nil
}

const AnywhereValueParameter = "..."

func ParameterNames() map[string]string {
	return map[string]string{
		"author_username":  "Имя автора",
		"from_branch_name": "Изначальная ветка",
		"to_branch_name":   "Конечная ветка",
		"status":           "Статус",
		"is_merge":         "Это мерж",
	}
}

func AllowParameters() map[string]map[string]map[string]string {
	return map[string]map[string]map[string]string{
		"pipeline": {
			"author_username": map[string]string{
				AnywhereValueParameter: "",
			},
			"from_branch_name": map[string]string{
				"develop":              "Develop",
				"master":               "Master",
				"release":              "Release",
				AnywhereValueParameter: "",
			},
			"to_branch_name": map[string]string{
				"develop":              "Develop",
				"master":               "Master",
				"release":              "Release",
				AnywhereValueParameter: "",
			},
			"status": map[string]string{
				"failed":  "Завершён с ошибкой",
				"success": "Успешно завершён",
			},
			"is_merge": map[string]string{
				"true":  "Да",
				"false": "Нет",
			},
		},
	}
}
