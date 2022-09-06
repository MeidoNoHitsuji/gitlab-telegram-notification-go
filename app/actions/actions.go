package actions

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab-telegram-notification-go/database"
	"gitlab-telegram-notification-go/telegram"
)

type ActionErrorType error

type ErrorForUser struct {
	ActionErrorType
}

func NewErrorForUser(s string) *ErrorForUser {
	return &ErrorForUser{
		ActionErrorType: errors.New(s),
	}
}

func GetActualActions() []BaseInterface {
	return []BaseInterface{
		NewBackTextAction(),
		NewBackCallbackAction(),
		NewTestAction(),
		NewStartAction(),
		NewSayAction(),
		NewTomatoFailAction(),
		NewSubscribesAction(),
		NewSelectProjectAction(),
		NewSelectProjectSettingsAction(),

		NewSubscribeAction(),
	}
}

func Active(update tgbotapi.Update) ActionErrorType {
	for _, action := range GetActualActions() {
		if action.Validate(update) {

			//TODO: Решить проблему с передачей CallbackData в Active
			err := action.Active(update)

			if err != nil {
				return err
			}

			UpdateActualAction(update, action)
			return nil
		}
	}

	return NewErrorForUser("Я не понимаю, что ты от меня хочешь.")
}

func UpdateActualAction(update tgbotapi.Update, action BaseInterface) {
	actionName := action.GetActionName()

	var chatId int64
	var username string

	if actionName != "" {

		if update.CallbackQuery != nil {
			chatId = update.CallbackQuery.Message.Chat.ID
			username = update.CallbackQuery.From.UserName
		} else if update.Message != nil {
			chatId = update.Message.Chat.ID
			username = update.Message.From.UserName
		} else {
			return
		}

		database.UpdateUserActionInChannel(
			chatId,
			username,
			string(actionName),
		)
	}
}

func GetActualAction(update tgbotapi.Update) ActionNameType {
	var message *tgbotapi.Message
	if update.CallbackQuery != nil {
		message = update.CallbackQuery.Message
	} else if update.Message != nil {
		message = update.Message
	}

	if message != nil {
		return ActionNameType(database.GetUserActionInChannel(
			message.Chat.ID,
			message.From.UserName,
		))
	}

	return ""
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
		beforeAction.SetIsBack()
		err := beforeAction.Active(update)
		if err != nil {
			return err
		}

		UpdateActualAction(update, beforeAction)
	}

	return nil
}
