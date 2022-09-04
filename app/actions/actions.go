package actions

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab-telegram-notification-go/database"
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
		NewStartAction(),
		NewSayAction(),
		NewTomatoFailAction(),
		NewSubscribesAction(),
		NewSelectProjectAction(),

		NewSubscribeAction(),
	}
}

func Active(update tgbotapi.Update) ActionErrorType {
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

	return NewErrorForUser("Я не понимаю, что ты от меня хочешь.")
}

func UpdateActualAction(update tgbotapi.Update, action BaseInterface) {
	a := action.GetActionName()
	var message *tgbotapi.Message

	if a != "" {
		if update.CallbackQuery != nil {
			message = update.CallbackQuery.Message
		} else if update.Message != nil {
			message = update.Message
		}

		if message != nil {
			database.UpdateUserActionInChannel(
				message.Chat.ID,
				message.From.UserName,
				string(a),
			)
		}
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

//TODO: Сносить старый кейборд
func BackAction(update tgbotapi.Update) error {
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
		err := beforeAction.Active(update)
		if err != nil {
			return err
		}

		//TODO: Добавить отдельную функцию у на удаления кейборда

		UpdateActualAction(update, beforeAction)
	}

	return nil
}
