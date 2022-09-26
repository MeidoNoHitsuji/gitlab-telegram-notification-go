package toggl

import (
	"errors"
	"fmt"
	"gitlab-telegram-notification-go/database"
	"gitlab-telegram-notification-go/models"
	"os"
)

func ActiveSubscription(telegramChannelId int64, active bool) error {
	var token models.UserToken

	db := database.Instant()

	var user models.User

	db.Where(models.User{
		TelegramChannelId: telegramChannelId,
	}).First(&user)

	res := db.Where(models.UserToken{
		TokenType: models.ToggleToken,
		UserId:    user.ID,
	}).First(&token)

	if res.RowsAffected == 0 {
		return errors.New("Токен toggle не был найден")
	}

	userData, err := Me(token.Token)

	if err != nil {
		return err
	}

	subscriptions, err := GetSubscriptions(userData.DefaultWorkspaceId, token.Token)

	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/%d", os.Getenv("TOGGLE_WEBHOOK_URL"), userData.Id)

	var needSubscription *SubscriptionData

	for _, subscription := range subscriptions {
		if subscription.UrlCallback == url {
			needSubscription = subscription
		}
	}

	if needSubscription != nil {
		_, err = EnableSubscription(userData.DefaultWorkspaceId, needSubscription.SubscriptionId, active, token.Token)
		return err
	} else {
		_, err = CreateSubscription(userData.DefaultWorkspaceId, token.Token, SubscriptionData{
			Description: "AtwintaTelegramBot subscription for Jira",
			Enabled:     active,
			UrlCallback: url,
			EventFilters: []SubscriptionEventData{
				{
					Action: "*",
					Entity: "time_entry",
				},
			},
		})

		return err
	}
}

func GetStatusSubscription(telegramChannelId int64) (bool, error) {
	var token models.UserToken
	var user models.User

	db := database.Instant()

	db.Where(models.User{
		TelegramChannelId: telegramChannelId,
	}).First(&user)

	res := db.Where(models.UserToken{
		TokenType: models.ToggleToken,
		UserId:    user.ID,
	}).First(&token)

	if res.RowsAffected == 0 {
		return false, errors.New("Токен toggle не был найден")
	}

	userData, err := Me(token.Token)

	if err != nil {
		return false, err
	}

	subscriptions, err := GetSubscriptions(userData.DefaultWorkspaceId, token.Token)

	url := fmt.Sprintf("%s/%d", os.Getenv("TOGGLE_WEBHOOK_URL"), userData.Id)

	if err != nil {
		return false, err
	}

	if len(subscriptions) == 0 {
		return false, nil
	} else {
		for _, subscription := range subscriptions {
			if subscription.UrlCallback == url {
				return subscription.Enabled, nil
			}
		}

		return false, nil
	}
}
