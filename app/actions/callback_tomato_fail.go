package actions

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab-telegram-notification-go/actions/callbacks"
	"gitlab-telegram-notification-go/telegram"
)

type TomatoFailAction struct {
	BaseAction
	CallbackData *callbacks.TomatoFailType
}

func (act *TomatoFailAction) Active(update tgbotapi.Update) error {
	act.CallbackData.Count++
	out, err := json.Marshal(act.CallbackData)

	if err != nil {
		return err
	}

	_, err = telegram.UpdateMessageById(
		update.CallbackQuery.Message,
		update.CallbackQuery.Message.Text,
		tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(
					fmt.Sprintf("%d 🍅", act.CallbackData.Count),
					string(out),
				),
			),
		),
		update.CallbackQuery.Message.Entities,
	)

	if err != nil {
		return err
	}

	return nil
}

func NewTomatoFailAction() *TomatoFailAction {
	return &TomatoFailAction{
		BaseAction: BaseAction{
			InitBy: InitByCallback,
		},
	}
}
