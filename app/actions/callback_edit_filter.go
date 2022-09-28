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
	fm "gitlab-telegram-notification-go/helper/formater"
	"gitlab-telegram-notification-go/models"
	"gitlab-telegram-notification-go/telegram"
	"strings"
)

const EditFilter ActionNameType = "ef_act" //edit_filter

type EditFilterActon struct {
	BaseAction
	CallbackData *callbacks.EditFilterType `json:"callback_data"`
}

type EditFilterBackData struct {
	ProjectId int `json:"pi"`
}

func (act *EditFilterActon) SetIsBack(update tgbotapi.Update) error {

	err := act.BaseAction.SetIsBack(update)

	if err != nil {
		return err
	}

	chatId, username := GetChatIdAndUsernameByUpdate(update)

	if chatId == 0 {
		return errors.New("Не найден пользователь для update.")
	}

	action := database.GetUserActionInChannel(chatId, username)

	if action == nil {
		return errors.New("Не найден action для update.")
	}

	if err = json.Unmarshal([]byte(action.Parameters), &act.CallbackData); err != nil {
		return err
	}

	UpdateActualActionParameter(update, "")

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
			AfterAction: []ActionNameType{
				CreateFilter,
				EditFilterParameter,
				SelectFilter,
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

		chatId, username := GetChatIdAndUsernameByUpdate(update)

		if chatId == 0 {
			return false
		}

		action := database.GetUserActionInChannel(chatId, username)

		if action == nil {
			return false
		}

		result := json.Unmarshal([]byte(action.Parameters), &act.CallbackData) == nil

		if !result {
			return false
		}

		act.CallbackData.ParameterValue = message.Text

		UpdateActualActionParameter(update, "")
		telegram.SendRemoveKeyboard(message.Chat.ID, false)
		return true
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

	subscribeObj := database.FirstOrCreateSubscribe(project.ID, message.Chat.ID, true)

	var subscribeEvent models.SubscribeEvent

	if act.CallbackData.EventId == 0 {
		subscribeEvent = models.SubscribeEvent{
			SubscribeId: subscribeObj.ID,
			Event:       act.CallbackData.EventName,
			Parameters:  map[string][]string{},
		}

		db.Save(&subscribeEvent)

		gitclient.SubscribeByProject(project)
	} else {
		subscribeEvent = models.SubscribeEvent{
			ID: act.CallbackData.EventId,
		}

		result := db.Find(&subscribeEvent)
		if subscribeEvent.Parameters == nil {
			subscribeEvent.Parameters = map[string][]string{}
		}

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

	// Редактируем данные, если это сейчас необходимо и передано!!
	if act.CallbackData.EditFormatter {
		if act.CallbackData.DeleteValue {
			subscribeEvent.Formatter = ""
			db.Save(&subscribeEvent)
		} else if act.CallbackData.FormatterValue != "" {
			subscribeEvent.Formatter = act.CallbackData.FormatterValue
			db.Save(&subscribeEvent)
		}
	} else if act.CallbackData.ParameterName != "" {
		if act.CallbackData.DeleteValue {
			subscribeEvent.Parameters[lowName[act.CallbackData.ParameterName]] = []string{}
			db.Save(&subscribeEvent)
		} else if act.CallbackData.ParameterValue != "" {
			ps := subscribeEvent.Parameters[lowName[act.CallbackData.ParameterName]]

			parameterValues := strings.Split(act.CallbackData.ParameterValue, ", ")

			for _, v := range parameterValues {
				if helper.Contains(ps, v) {
					ps = helper.Drop(ps, v)
				} else {
					ps = append(ps, v)
				}
			}

			subscribeEvent.Parameters[lowName[act.CallbackData.ParameterName]] = ps
			db.Save(&subscribeEvent)
		}
	}
	//Закончили редактировать данные!!

	var text string
	var keyboards [][]tgbotapi.InlineKeyboardButton
	var parameterKeys []string

	allowParameters := AllowParameters()
	parameterNames := ParameterNames()
	parameters, canEditParameters := allowParameters[subscribeEvent.Event]

	//Вытаскиваем только те ключи, которые мы можем обрабатывать через webhook'и
	if canEditParameters {
		for k := range parameters {
			parameterKeys = append(parameterKeys, k)
		}
	}

	allowFormatters := AllowFormatters()
	formatters, canEditFormatter := allowFormatters[subscribeEvent.Event]

	//Формируем блок текста для параметров
	if canEditParameters {

		allowEvents := helper.AllowEventsWithName()

		text = fmt.Sprintf("Тип ивента: %s\n—————\nИмеющиеся параметры:", allowEvents[subscribeEvent.Event])

		if len(subscribeEvent.Parameters) != 0 {
			for _, key := range parameterKeys {
				currentParams, ok := subscribeEvent.Parameters[lowName[key]]

				if ok && len(currentParams) > 0 {
					var newCurrentParams []string

					for _, param := range currentParams {
						newCurrentParams = append(newCurrentParams, fm.Underline(param))
					}

					text = fmt.Sprintf("%s\n> %s: %s", text, parameterNames[key], strings.Join(newCurrentParams, ", "))
				} else {
					text = fmt.Sprintf("%s\n> %s: %s", text, parameterNames[key], fm.Italic("Ну учитывается"))
				}
			}
		} else {
			text = fmt.Sprintf("%s\nОтсутствуют.", text)
		}

		if act.CallbackData.ParameterName != "" {
			text = fmt.Sprintf("%s\n\nРедактируемый параметр: %s", text, parameterNames[act.CallbackData.ParameterName])
		}
	} else {
		text = "Фильтры данного ивента ещё нельзя редактировать!"
	}
	//Закончили фромировать блок текста для параметров

	//Формируем блок текста для форматтера
	if canEditFormatter {
		currentFormatter, ok := formatters[subscribeEvent.Formatter]

		if ok {
			text = fmt.Sprintf("%s\n—————\nТип форматирования: %s", text, currentFormatter)
		} else {
			text = fmt.Sprintf("%s\n—————\nТип форматирования: По умолчанию", text)
		}
	} else {
		text = fmt.Sprintf("%s\n—————\nУ данного ивента ещё нельзя менять форматирование", text)
	}
	//Закончили фромировать блок текста для

	//Если мы сейчас редактируем форматтер
	if act.CallbackData.EditFormatter {
		//Добавляем кнопки форматтера
		if canEditFormatter {
			for _, keys := range helper.Grouping(helper.Keys(formatters), 3) {
				var keyboardButtons []tgbotapi.InlineKeyboardButton
				for j := 0; j < len(keys); j++ {

					eventData := callbacks.NewEditFilterWithFormatterValueType(project.ID, subscribeEvent.ID, keys[j])
					eventOut, err := json.Marshal(eventData)

					if err != nil {
						return err
					}

					keyboardButtons = append(
						keyboardButtons,
						tgbotapi.NewInlineKeyboardButtonData(formatters[keys[j]], string(eventOut)),
					)
				}

				keyboards = append(keyboards, tgbotapi.NewInlineKeyboardRow(keyboardButtons...))
			}
		}

		deleteFormatterOut, err := json.Marshal(
			callbacks.NewEditFilterWithDeleteFormatterType(project.ID, subscribeEvent.ID),
		)

		if err != nil {
			return err
		}

		backOut, err := json.Marshal(
			callbacks.NewEditFilterWithEventIdType(project.ID, subscribeEvent.ID),
		)

		if err != nil {
			return err
		}

		keyboards = append(keyboards,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("По умолчанию", string(deleteFormatterOut)),
				tgbotapi.NewInlineKeyboardButtonData("Назад", string(backOut)),
			),
		)
		//Иначе, если мы редактируем параметры
	} else if act.CallbackData.ParameterName != "" {
		//И можем это делать
		if canEditParameters {
			//То добавляем кнопки параметров, которые нам доступны
			filters := parameters[act.CallbackData.ParameterName]
			filterKeys := helper.Keys(filters)
			if helper.Contains(filterKeys, AnywhereValueParameter) {

				eventData := callbacks.NewEditFilterWithParameterType(project.ID, subscribeEvent.ID, act.CallbackData.ParameterName)
				eventOut, err := json.Marshal(eventData)

				if err != nil {
					return err
				}
				UpdateActualActionParameter(update, string(eventOut))

				buttonOut, err := json.Marshal(
					callbacks.NewEditFilterParameterType(project.ID, subscribeEvent.ID, act.CallbackData.ParameterName),
				)

				if err != nil {
					return err
				}

				keyboards = append(keyboards,
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("Ввести значение", string(buttonOut)),
					),
				)
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
	}

	//Добавляем дефолтные кнопки из главного меню (если мы не редактируем ни то, ни то)
	if act.CallbackData.ParameterName == "" && !act.CallbackData.EditFormatter {

		// Если мы можем редактировать параметры
		if canEditParameters {
			//То добавляем кнопки возможных параметров для редактирования
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
		}

		//Если мы можем редактировать форматтер
		if canEditFormatter && len(formatters) > 0 {
			formatterOut, err := json.Marshal(
				callbacks.NewEditFilterWithFormatterType(project.ID, subscribeEvent.ID),
			)

			if err != nil {
				return err
			}

			keyboardFormatter := tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Изменить форматирование", string(formatterOut)),
			)

			keyboards = append(keyboards, keyboardFormatter)
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
	//Закончили добавлять дефолтные кнопки

	keyboards = append(keyboards, keyboardBack)

	if act.IsBack || act.InitializedBy == InitByText {
		telegram.SendMessageById(
			message.Chat.ID,
			text,
			tgbotapi.NewInlineKeyboardMarkup(keyboards...),
			nil,
		)
	} else {
		telegram.UpdateMessageWithParseById(
			message,
			text,
			tgbotapi.NewInlineKeyboardMarkup(keyboards...),
		)
	}

	return nil
}

const AnywhereValueParameter = "..."

const (
	AuthorUsernameParameter = "au"
	FromBranchNameParameter = "fbn"
	ToBranchNameParameter   = "tbn"
	StatusParameter         = "s"
	IsMergeParameter        = "im"
)

var lowName = map[string]string{
	AuthorUsernameParameter: "author_username",
	FromBranchNameParameter: "from_branch_name",
	ToBranchNameParameter:   "to_branch_name",
	StatusParameter:         "status",
	IsMergeParameter:        "is_merge",
}

func ParameterNames() map[string]string {
	return map[string]string{
		AuthorUsernameParameter: "Никнейм автора",
		FromBranchNameParameter: "Изначальная ветка",
		ToBranchNameParameter:   "Конечная ветка",
		StatusParameter:         "Статус",
		IsMergeParameter:        "Это мерж",
	}
}

func AllowParameters() map[string]map[string]map[string]string {
	return map[string]map[string]map[string]string{
		"pipeline": {
			AuthorUsernameParameter: map[string]string{
				AnywhereValueParameter: "",
			},
			FromBranchNameParameter: map[string]string{
				"develop":              "",
				"master":               "",
				"release":              "",
				AnywhereValueParameter: "",
			},
			ToBranchNameParameter: map[string]string{
				"develop":              "",
				"master":               "",
				"release":              "",
				AnywhereValueParameter: "",
			},
			StatusParameter: map[string]string{
				"failed":  "Завершён с ошибкой",
				"success": "Успешно завершён",
			},
			IsMergeParameter: map[string]string{
				"true":  "Да",
				"false": "Нет",
			},
		},
	}
}

func AllowFormatters() map[string]map[string]string {
	return map[string]map[string]string{
		"pipeline": {
			"default": "Пункты сборки",
			"commits": "Новые комиты",
			"logs":    "Соглашение о комитах с ссылками на Jira",
		},
	}
}
