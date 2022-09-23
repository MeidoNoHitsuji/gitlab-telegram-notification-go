package actions

import (
	"encoding/json"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab-telegram-notification-go/actions/callbacks"
	"gitlab-telegram-notification-go/actions/middlewares"
	"gitlab-telegram-notification-go/telegram"
)

const UserSettingEnterTokenActionType ActionNameType = "uset_act" //user_setting_enter_token

type UserSettingEnterTokenAction struct {
	BaseAction
	CallbackData *callbacks.UserSettingEnterTokenType `json:"callback_data"`
}

func NewUserSettingEnterTokenAction() *UserSettingEnterTokenAction {
	return &UserSettingEnterTokenAction{
		BaseAction: BaseAction{
			ID:                    UserSettingEnterTokenActionType,
			InitBy:                []ActionInitByType{InitByCallback},
			InitCallbackFuncNames: []callbacks.CallbackFuncName{callbacks.UserSettingEnterTokenFuncName},
			BeforeAction:          UserSettingTokensActionType,
			MiddleWares:           []middlewares.MiddleWares{middlewares.OnlyDM},
		},
	}
}

func (act *UserSettingEnterTokenAction) Validate(update tgbotapi.Update) bool {
	if !act.BaseAction.Validate(update) {
		return false
	}

	res := json.Unmarshal([]byte(update.CallbackQuery.Data), &act.CallbackData) == nil

	if !res {
		return false
	}

	eventData := callbacks.NewUserSettingTokensWithTokenType(act.CallbackData.TokenType)
	eventOut, err := json.Marshal(eventData)

	if err != nil {
		return false
	}

	UpdateActualActionParameter(update, string(eventOut))

	return true
}

func (act *UserSettingEnterTokenAction) Active(update tgbotapi.Update) error {
	message, _ := telegram.GetMessageFromUpdate(update)

	if message == nil {
		return errors.New("Неизвестно откуда прилетел запрос.")
	}

	var keyboards [][]tgbotapi.KeyboardButton

	keyboards = append(keyboards, tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Удалить значение"),
	), tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Отмена"),
	))

	keyboard := tgbotapi.NewReplyKeyboard(keyboards...)
	keyboard.OneTimeKeyboard = true

	actualIntegration := ActualIntegrationTokens()

	text := fmt.Sprintf("Введите значение токена для интеграцией с %s", actualIntegration[act.CallbackData.TokenType])

	telegram.SendMessageById(
		message.Chat.ID,
		text,
		keyboard,
		nil,
	)

	return nil
}
