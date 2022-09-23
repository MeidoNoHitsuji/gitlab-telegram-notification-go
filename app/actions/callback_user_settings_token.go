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
	"gitlab-telegram-notification-go/models"
	"gitlab-telegram-notification-go/telegram"
)

const UserSettingTokensActionType ActionNameType = "ust_act" //user_setting_tokens

type UserSettingTokensAction struct {
	BaseAction
	CallbackData *callbacks.UserSettingTokensType `json:"callback_data"`
}

func NewUserSettingTokensAction() *UserSettingTokensAction {
	return &UserSettingTokensAction{
		BaseAction: BaseAction{
			ID: UserSettingTokensActionType,
			InitBy: []ActionInitByType{
				InitByCallback,
				InitByText,
			},
			InitCallbackFuncNames: []callbacks.CallbackFuncName{callbacks.UserSettingTokensFuncName},
			AfterAction: []ActionNameType{
				UserSettingsActionType,
				UserSettingEnterTokenActionType,
			},
			BeforeAction: UserSettingsActionType,
			MiddleWares:  []middlewares.MiddleWares{middlewares.OnlyDM},
		},
	}
}

func (act *UserSettingTokensAction) Validate(update tgbotapi.Update) bool {
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

		act.CallbackData.TokenValue = message.Text

		UpdateActualActionParameter(update, "")
		telegram.SendRemoveKeyboard(message.Chat.ID, false)
	}

	return true
}

func (act *UserSettingTokensAction) Active(update tgbotapi.Update) error {
	message, _ := telegram.GetMessageFromUpdate(update)

	if message == nil {
		return errors.New("Неизвестно откуда прилетел запрос.")
	}

	actualIntegrations := ActualIntegrationTokens()
	db := database.Instant()
	chatId, _ := GetChatIdAndUsernameByUpdate(update)

	if act.CallbackData != nil {
		if act.CallbackData.TokenType != "" && act.CallbackData.TokenValue != "" {
			var token models.UserToken

			user := models.User{
				TelegramChannelId: chatId,
			}

			res := db.Where(models.UserToken{
				User:      user,
				TokenType: act.CallbackData.TokenType,
			}).First(&token)

			if act.CallbackData.TokenValue == "Удалить значение" {
				if res.RowsAffected != 0 {
					db.Unscoped().Delete(&token)
				}
			} else {
				token.Token = act.CallbackData.TokenValue
				if res.RowsAffected == 0 {
					db.First(&user)
					token.UserId = user.ID
					token.TokenType = act.CallbackData.TokenType
					db.Omit("User").Create(&token)
				} else {
					db.Omit("User").Save(&token)
				}
			}
		}
	}

	text := "Тут ты можешь поменять токены от различных сервисов, с которыми настроена интеграция."

	var tokens []models.UserToken

	db.Where(models.UserToken{
		User: models.User{
			TelegramChannelId: chatId,
		},
	}).Find(&tokens)

	text = fmt.Sprintf("%s\n—————\nИмеющиеся токены:", text)

	for tokenType, name := range actualIntegrations {
		tokenValue := "Отсутствует"
		for _, token := range tokens {
			if token.TokenType == tokenType {
				tokenValue = token.Token
			}
		}
		text = fmt.Sprintf("%s\n> %s: %s", text, name, tokenValue)
	}

	var keyboards [][]tgbotapi.InlineKeyboardButton

	for _, keys := range helper.Grouping(helper.Keys(actualIntegrations), 3) {
		var keyboardButtons []tgbotapi.InlineKeyboardButton
		for j := 0; j < len(keys); j++ {

			eventOut, err := json.Marshal(
				callbacks.NewUserSettingEnterTokenType(keys[j]),
			)

			if err != nil {
				return err
			}

			keyboardButtons = append(
				keyboardButtons,
				tgbotapi.NewInlineKeyboardButtonData(actualIntegrations[keys[j]], string(eventOut)),
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

	if act.InitializedBy == InitByText || act.IsBack {
		telegram.SendMessageById(
			message.Chat.ID,
			text,
			tgbotapi.NewInlineKeyboardMarkup(keyboards...),
			nil,
		)
	} else {
		telegram.UpdateMessageById(
			message,
			text,
			tgbotapi.NewInlineKeyboardMarkup(keyboards...),
			nil,
		)
	}

	return nil
}

func ActualIntegrationTokens() map[string]string {
	return map[string]string{
		"jira":   "Jira",
		"toggle": "Toggle",
		"gitlab": "Gitlab",
	}
}
