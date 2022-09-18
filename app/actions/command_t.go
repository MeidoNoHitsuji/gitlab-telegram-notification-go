package actions

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab-telegram-notification-go/telegram"
)

const TestActionType ActionNameType = "test"

type TestAction struct {
	BaseAction
}

func NewTestAction() *TestAction {
	return &TestAction{
		BaseAction: BaseAction{
			ID:       TestActionType,
			InitBy:   []ActionInitByType{InitByCommand},
			InitText: "test",
		},
	}
}

func (act *TestAction) Active(update tgbotapi.Update) error {
	//msg, err := telegram.SendMessageById(update.Message.Chat.ID, "remove_keyboard", tgbotapi.NewRemoveKeyboard(false), nil)
	//
	//if err != nil {
	//	return err
	//}
	//
	//bot := telegram.Instant()
	//
	//del := tgbotapi.NewDeleteMessage(msg.Chat.ID, msg.MessageID)
	//
	//bot.Request(del)

	cfg := tgbotapi.GetChatMenuButtonConfig{
		ChatID: update.Message.Chat.ID,
	}

	bot := telegram.Instant()

	data, err := bot.Request(cfg)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	fmt.Println(data)

	return nil
}
