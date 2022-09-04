package actions

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab-telegram-notification-go/actions/callbacks"
	"gitlab-telegram-notification-go/database"
	"gitlab-telegram-notification-go/telegram"
)

const SubscribesActionType ActionNameType = "subscribes"

type SubscribesAction struct {
	BaseAction
	CallbackData *callbacks.SubscribesType
}

func NewSubscribesAction() *SubscribesAction {
	return &SubscribesAction{
		BaseAction: BaseAction{
			ID:                   SubscribesActionType,
			InitBy:               InitByCallback,
			InitCallbackFuncName: callbacks.SubscribesFuncName,
			BeforeAction:         Start,
		},
	}
}

func (act *SubscribesAction) Active(update tgbotapi.Update) error {

	var message *tgbotapi.Message

	if update.Message != nil {
		message = update.Message
	} else if update.CallbackQuery != nil {
		message = update.CallbackQuery.Message
	} else {
		return errors.New("Неизвестно откуда прилетел запрос.")
	}

	projects := database.GetProjectsByTelegramIds(message.Chat.ID)

	var keyboard [][]tgbotapi.KeyboardButton
	lines := len(projects) / 3

	if len(projects)%3 > 0 {
		lines++
	}

	for i := 0; i < lines; i++ {
		pr := projects[i*3 : ((i + 1) * 3)]
		var keyboardButtons []tgbotapi.KeyboardButton
		for j := 0; j < len(pr); j++ {
			keyboardButtons = append(keyboardButtons, tgbotapi.NewKeyboardButton(pr[j].Name))
		}
		keyboard = append(keyboard, tgbotapi.NewKeyboardButtonRow(keyboardButtons...))
	}

	keyboard = append(keyboard, tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Отмена"),
	))

	telegram.SendMessageById(
		message.Chat.ID,
		`Чтобы добавить/обновить подписку вам необходимо выбрать проект для этого.
Если вы не находите необходимый вам проект в списке, то введите его slug или id.`,
		tgbotapi.NewReplyKeyboard(keyboard...),
		nil,
	)

	return nil
}
