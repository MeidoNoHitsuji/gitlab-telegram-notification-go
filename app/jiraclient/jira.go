package jiraclient

import (
	"errors"
	"github.com/andygrunwald/go-jira"
	"gitlab-telegram-notification-go/database"
	"gitlab-telegram-notification-go/models"
	"os"
)

func New(telegramChannelId int64) (*jira.Client, error) {
	var user models.User
	var userToken models.UserToken

	db := database.Instant()

	res := db.Where(models.User{
		TelegramChannelId: telegramChannelId,
	}).Find(&user)

	if res.RowsAffected == 0 {
		return nil, errors.New("Пользователь не был найден по telegramChannelId")
	}

	res = db.Where(models.UserToken{
		UserId:    user.ID,
		TokenType: models.JiraToken,
	}).Find(&userToken)

	if res.RowsAffected == 0 {
		return nil, errors.New("JiraToken не был найден у пользователя")
	}

	patTransport := &jira.PATAuthTransport{
		Token: userToken.Token,
	}

	return jira.NewClient(patTransport.Client(), os.Getenv("JIRA_DOMAIN"))
}
