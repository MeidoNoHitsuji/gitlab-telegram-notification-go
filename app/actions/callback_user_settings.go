package actions

import (
	"encoding/json"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab-telegram-notification-go/actions/callbacks"
	"gitlab-telegram-notification-go/actions/middlewares"
	"gitlab-telegram-notification-go/telegram"
)

const UserSettingsActionType ActionNameType = "us_act" //user_settings

type UserSettingsAction struct {
	BaseAction
}

func NewUserSettingsAction() *UserSettingsAction {
	return &UserSettingsAction{
		BaseAction: BaseAction{
			ID:                    UserSettingsActionType,
			InitBy:                []ActionInitByType{InitByCallback},
			InitCallbackFuncNames: []callbacks.CallbackFuncName{callbacks.UserSettingsFuncName},
			BeforeAction:          Start,
			MiddleWares:           []middlewares.MiddleWares{middlewares.OnlyDM},
		},
	}
}

func (act *UserSettingsAction) Validate(update tgbotapi.Update) bool {
	return act.BaseAction.Validate(update)
}

func (act *UserSettingsAction) Active(update tgbotapi.Update) error {
	message, _ := telegram.GetMessageFromUpdate(update)

	if message == nil {
		return errors.New("Неизвестно откуда прилетел запрос.")
	}

	var keyboards [][]tgbotapi.InlineKeyboardButton

	tokensOut, err := json.Marshal(
		callbacks.NewUserSettingTokensType(),
	)

	if err != nil {
		return err
	}

	integrationsOut, err := json.Marshal(
		callbacks.NewUserIntegrationsType(),
	)

	if err != nil {
		return err
	}

	keyboards = append(keyboards, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Интеграции", string(integrationsOut)),
		tgbotapi.NewInlineKeyboardButtonData("Токены", string(tokensOut)),
	))

	backOut, err := json.Marshal(
		callbacks.NewBackType(
			callbacks.NewStartType(),
		),
	)

	if err != nil {
		return err
	}

	keyboards = append(keyboards, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Отмена", string(backOut)),
	))

	telegram.UpdateMessageById(
		message,
		"Выбери какие настройки пользователя ты хочешь изменить.",
		tgbotapi.NewInlineKeyboardMarkup(keyboards...),
		nil,
	)

	return nil
}
