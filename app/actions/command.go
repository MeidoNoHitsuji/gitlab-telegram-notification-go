package actions

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab-telegram-notification-go/database"
	"gitlab-telegram-notification-go/telegram"
)

func Test(telegramId ...int64) {
	senderId := telegramId[0]

	projects := database.GetProjectsByTelegramIds(telegramId...)

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

	telegram.SendMessageById(senderId, "Это какая-то хуита?", tgbotapi.NewReplyKeyboard(keyboard...), nil)
}
