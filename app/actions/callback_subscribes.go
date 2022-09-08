package actions

import (
	"encoding/json"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab-telegram-notification-go/actions/callbacks"
	"gitlab-telegram-notification-go/database"
	"gitlab-telegram-notification-go/telegram"
)

const SubscribesActionType ActionNameType = "subscribes"

type SubscribesAction struct {
	BaseAction
	CallbackData *callbacks.SubscribesType `json:"callback_data"`
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

func (act *SubscribesAction) Validate(update tgbotapi.Update) bool {
	if !act.BaseAction.Validate(update) {
		return false

	}
	return json.Unmarshal([]byte(update.CallbackQuery.Data), &act.CallbackData) == nil
}

func (act *SubscribesAction) Active(update tgbotapi.Update) error {
	message, _ := telegram.GetMessageFromUpdate(update)

	if message == nil {
		return errors.New("Неизвестно откуда прилетел запрос.")
	}

	projects := database.GetProjectsByTelegramIds(message.Chat.ID)

	var keyboardRows [][]tgbotapi.KeyboardButton
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
		keyboardRows = append(keyboardRows, tgbotapi.NewKeyboardButtonRow(keyboardButtons...))
	}

	keyboardRows = append(keyboardRows, tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Отмена"),
	))

	keyboard := tgbotapi.NewReplyKeyboard(keyboardRows...)
	keyboard.OneTimeKeyboard = true

	telegram.SendMessageById(
		message.Chat.ID,
		"Чтобы добавить/обновить подписку вам необходимо выбрать проект для этого\nЕсли вы не находите необходимый вам проект в списке, то введите его slug или id.",
		keyboard,
		nil,
	)

	return nil
}
