package actions

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab-telegram-notification-go/database"
	"gitlab-telegram-notification-go/telegram"
	"strings"
)

type ErrorForUser struct {
	Message string
}

func (e ErrorForUser) Error() string {
	return e.Message
}

func NewErrorForUser(s string) ErrorForUser {
	return ErrorForUser{
		Message: s,
	}
}

func GetActualActions() []BaseInterface {
	return []BaseInterface{
		NewTomatoFailAction(),

		NewBackTextAction(),
		NewBackCallbackAction(),
		NewTestAction(),
		NewStartAction(),
		NewUserSettingsAction(),
		NewUserIntegrationsAction(),
		NewUserSettingTokensAction(),
		NewUserSettingEnterTokenAction(),
		NewSayAction(),
		NewStopKeyboardAction(),
		NewSubscribesAction(),
		NewSelectProjectAction(),
		NewSelectProjectSettingsAction(),
		NewCreateFilterActon(),
		NewSelectFilterActon(),
		NewEditFilterActon(),
		NewEditFilterParameterActon(),

		NewSubscribeAction(),
	}
}

func Active(update tgbotapi.Update) error {
	for _, action := range GetActualActions() {
		if action.Validate(update) {

			err := action.Active(update)

			if err != nil {
				return err
			}

			UpdateActualAction(update, action)
			return nil
		}
	}

	return errors.New("Я не понимаю, что от меня хочет пользователь.")
}

func UpdateActualActionParameter(update tgbotapi.Update, parameter string) {
	chatId, username := GetChatIdAndUsernameByUpdate(update)

	if chatId == 0 {
		return
	}

	database.UpdateUserActionParameterInChannel(
		chatId,
		username,
		parameter,
	)
}

func UpdateActualAction(update tgbotapi.Update, action BaseInterface) {
	actionName := action.GetActionName()

	if actionName != "" {

		chatId, username := GetChatIdAndUsernameByUpdate(update)

		if chatId == 0 {
			return
		}

		database.UpdateUserActionInChannel(
			chatId,
			strings.ToLower(username),
			string(actionName),
		)
	}
}

func GetChatIdAndUsernameByUpdate(update tgbotapi.Update) (int64, string) {
	if update.CallbackQuery != nil {
		return update.CallbackQuery.Message.Chat.ID, strings.ToLower(update.CallbackQuery.From.UserName)
	} else if update.Message != nil {
		return update.Message.Chat.ID, strings.ToLower(update.Message.From.UserName)
	} else {
		return 0, ""
	}
}

func GetActualAction(update tgbotapi.Update) ActionNameType {
	chatId, username := GetChatIdAndUsernameByUpdate(update)

	if chatId == 0 {
		return ""
	}

	action := database.GetUserActionInChannel(chatId, username)

	if action != nil {
		return ActionNameType(action.Action)
	} else {
		return ""
	}
}

func BackAction(update tgbotapi.Update) error {
	message, _ := telegram.GetMessageFromUpdate(update)
	actions := GetActualActions()
	actualAction := GetActualAction(update)
	var action BaseInterface

	for _, act := range actions {
		if action != nil {
			continue
		}
		if act.GetActionName() == actualAction {
			action = act
		}
	}

	if action == nil {
		return nil
	}

	beforeActionName := action.GetBeforeAction()

	if beforeActionName == "" {
		return nil
	}

	var beforeAction BaseInterface

	for _, act := range actions {
		if beforeAction != nil {
			continue
		}
		if act.GetActionName() == beforeActionName {
			beforeAction = act
		}
	}

	if beforeAction != nil {
		telegram.SendRemoveKeyboard(message.Chat.ID, false)
		err := beforeAction.SetIsBack(update)
		if err != nil {
			return err
		}

		err = beforeAction.Active(update)
		if err != nil {
			return err
		}

		UpdateActualAction(update, beforeAction)
	}

	return nil
}
