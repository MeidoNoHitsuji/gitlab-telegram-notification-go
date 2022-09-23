package actions

import (
	"encoding/json"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab-telegram-notification-go/actions/callbacks"
	"gitlab-telegram-notification-go/database"
	"gitlab-telegram-notification-go/models"
	"gitlab-telegram-notification-go/telegram"
)

const Start ActionNameType = "start"

type StartAction struct {
	BaseAction
}

func NewStartAction() *StartAction {
	return &StartAction{
		BaseAction: BaseAction{
			ID: Start,
			InitBy: []ActionInitByType{
				InitByCommand,
				InitByCallback,
			},
			InitCallbackFuncNames: []callbacks.CallbackFuncName{callbacks.StartFuncName},
			InitText:              "start",
		},
	}
}

func (act *StartAction) Validate(update tgbotapi.Update) bool {
	return act.BaseAction.Validate(update)
}

func (act *StartAction) Active(update tgbotapi.Update) error {

	message, botMessage := telegram.GetMessageFromUpdate(update)

	if message == nil {
		return errors.New("Неизвестно откуда прилетел запрос.")
	}

	var keyboard []tgbotapi.InlineKeyboardButton

	subscribeOut, err := json.Marshal(
		callbacks.NewSubscribesType(),
	)

	if err != nil {
		return err
	}

	settingsOut, err := json.Marshal(
		callbacks.NewUserSettingsType(),
	)

	if err != nil {
		return err
	}

	db := database.Instant()
	chatId, _ := GetChatIdAndUsernameByUpdate(update)

	var user models.User

	result := db.Where(models.User{
		TelegramChannelId: chatId,
	}).First(&user)

	if result.RowsAffected != 0 {
		//Мы в личной переписке
		keyboard = tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Подписки", string(subscribeOut)),
			tgbotapi.NewInlineKeyboardButtonData("Настройки", string(settingsOut)),
		)
	} else {
		//Мы в чате
		keyboard = tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Подписки", string(subscribeOut)),
		)
	}

	if botMessage {
		telegram.UpdateMessageById(
			message,
			"Привет! Что поделаем?",
			tgbotapi.NewInlineKeyboardMarkup(keyboard),
			nil,
		)
	} else {
		telegram.SendMessageById(
			message.Chat.ID,
			"Привет! Что поделаем?",
			tgbotapi.NewInlineKeyboardMarkup(keyboard),
			nil,
		)
	}

	return nil
}
