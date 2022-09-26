package actions

import (
	"encoding/json"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab-telegram-notification-go/actions/callbacks"
	"gitlab-telegram-notification-go/actions/middlewares"
	"gitlab-telegram-notification-go/database"
	"gitlab-telegram-notification-go/helper"
	fm "gitlab-telegram-notification-go/helper/formater"
	"gitlab-telegram-notification-go/models"
	"gitlab-telegram-notification-go/telegram"
)

const UserIntegrationsActionType ActionNameType = "ui_act" //user_integrations

type UserIntegrationsAction struct {
	BaseAction
	CallbackData *callbacks.UserIntegrationsType `json:"callback_data"`
}

func NewUserIntegrationsAction() *UserIntegrationsAction {
	return &UserIntegrationsAction{
		BaseAction: BaseAction{
			ID:     UserIntegrationsActionType,
			InitBy: []ActionInitByType{InitByCallback},
			InitCallbackFuncNames: []callbacks.CallbackFuncName{
				callbacks.UserIntegrationsFuncName,
			},
			AfterAction:  []ActionNameType{UserSettingsActionType},
			BeforeAction: UserSettingsActionType,
			MiddleWares:  []middlewares.MiddleWares{middlewares.OnlyDM},
		},
	}
}

func (act *UserIntegrationsAction) Validate(update tgbotapi.Update) bool {
	if !act.BaseAction.Validate(update) {
		return false

	}
	return json.Unmarshal([]byte(update.CallbackQuery.Data), &act.CallbackData) == nil
}

func (act *UserIntegrationsAction) Active(update tgbotapi.Update) error {
	message, _ := telegram.GetMessageFromUpdate(update)

	if message == nil {
		return errors.New("Неизвестно откуда прилетел запрос.")
	}

	var token models.UserToken

	db := database.Instant()
	chatId, _ := GetChatIdAndUsernameByUpdate(update)

	switch act.CallbackData.IntegrationType {
	case models.ToggleJiraIntegration:
		res := db.Where(models.UserToken{
			User: models.User{
				TelegramChannelId: chatId,
			},
			TokenType: models.ToggleToken,
		}).First(&token)

		if res.RowsAffected == 0 || token.Token == "" {
			return NewErrorForUser("У вас отсутствует токен Toggle")
		}

		res = db.Where(models.UserToken{
			User: models.User{
				TelegramChannelId: chatId,
			},
			TokenType: models.JiraToken,
		}).First(&token)

		if res.RowsAffected == 0 || token.Token == "" {
			return NewErrorForUser("У вас отсутствует токен Jira")
		}

		//var integration models.UserIntegrations
		//
		//res = db.Where(models.UserIntegrations{
		//	IntegrationType: models.ToggleJiraIntegration,
		//	User: models.User{
		//		TelegramChannelId: chatId,
		//	},
		//}).First(&integration)
		//
		//if res.RowsAffected == 0 {
		//	integration.Active = true
		//	db.Omit("User").Create(&integration)
		//} else {
		//	integration.Active = !integration.Active
		//	db.Omit("User").Save(&integration)
		//}

		break
	default:
		break
	}

	var keyboards [][]tgbotapi.InlineKeyboardButton
	text := "Тут вы можете включить/отключить различные интеграции.\n\nСписок интеграция:"

	allowIntegrations := AllowIntegrations()

	i := 0

	var keysIntegration []string
	for key, data := range allowIntegrations {
		i++
		keysIntegration = append(keysIntegration, key)
		var integration models.UserIntegrations
		var t string
		result := db.Where(models.UserIntegrations{
			User: models.User{
				TelegramChannelId: chatId,
			},
			IntegrationType: key,
		}).First(&integration)

		if result.RowsAffected == 0 || !integration.Active {
			t = "Выкл"
		} else {
			t = "Вкл"
		}

		text = fmt.Sprintf("%s\n%d. %s: %s.\n%s", text, i, fm.Underline(data["title"]), fm.Bold(t), fm.Italic(data["description"]))
	}

	for _, keys := range helper.Grouping(keysIntegration, 3) {
		var keyboardButtons []tgbotapi.InlineKeyboardButton
		for j := 0; j < len(keys); j++ {

			eventOut, err := json.Marshal(
				callbacks.NewUserIntegrationsWithTypeType(keys[j]),
			)

			if err != nil {
				return err
			}

			keyboardButtons = append(
				keyboardButtons,
				tgbotapi.NewInlineKeyboardButtonData(allowIntegrations[keys[j]]["title"], string(eventOut)),
			)
		}

		keyboards = append(keyboards, tgbotapi.NewInlineKeyboardRow(keyboardButtons...))
	}

	backOut, err := json.Marshal(
		callbacks.NewBackType(
			callbacks.NewUserSettingsType(),
		),
	)

	if err != nil {
		return err
	}

	keyboards = append(keyboards, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Отмена", string(backOut)),
	))

	telegram.UpdateMessageWithParseById(
		message,
		text,
		tgbotapi.NewInlineKeyboardMarkup(keyboards...),
	)

	return nil
}

func AllowIntegrations() map[string]map[string]string {
	return map[string]map[string]string{
		models.ToggleJiraIntegration: {
			"title":       "Toggle - Jira",
			"description": "Интеграция времени между Toggle и Jira. Если у вашего трекера в Toggle будет в начале номер карточки, то время, которое вы трекали, будет автоматически добавлено в соответствующую карточку, если у вас есть доступ до неё. Для этого нужно обязательно заполнить токены Toggle и Jira.",
		},
	}
}
