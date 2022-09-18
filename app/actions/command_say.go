package actions

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab-telegram-notification-go/telegram"
)

const Say ActionNameType = "say"

type SayAction struct {
	BaseAction
}

func NewSayAction() *SayAction {
	return &SayAction{
		BaseAction: BaseAction{
			ID:       Say,
			InitBy:   []ActionInitByType{InitByCommand},
			InitText: "say",
		},
	}
}

func (act *SayAction) Active(update tgbotapi.Update) error {
	deleteMessageConfig := tgbotapi.NewDeleteMessage(update.Message.Chat.ID, update.Message.MessageID)
	bot := telegram.Instant()
	_, err := bot.Request(deleteMessageConfig)

	if err != nil {
		fmt.Println(err)
	}

	telegram.SendMessageById(update.Message.Chat.ID, update.Message.CommandArguments(), nil, update.Message.Entities)
	return nil
}
