package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab-telegram-notification-go/database"
	"gitlab-telegram-notification-go/models"
	"gitlab-telegram-notification-go/toggl"
	"os"
)

var debugCmd = &cobra.Command{
	Use: "debug",
	Run: debug,
}

func init() {
	rootCmd.AddCommand(debugCmd)
}

func debug(cmd *cobra.Command, args []string) {

	var token models.UserToken

	db := database.Instant()

	db.Where(models.UserToken{
		TokenType: "toggle",
		User: models.User{
			TelegramChannelId: 479413765,
		},
	}).First(&token)

	userData, err := toggl.Me(token.Token)

	if err != nil {
		panic(err)
	}

	result, err := toggl.Events()

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(result)
	}

	err = toggl.GetSubscriptions(userData.DefaultWorkspaceId, token.Token)

	url := fmt.Sprintf("%s/%s/%s/%d", os.Getenv("WEBHOOK_DOMAIN"), os.Getenv("WEBHOOK_URL"), os.Getenv("TOGGLE_WEBHOOK_URL"), userData.Id)

	err = toggl.CreateSubscriptions(userData.DefaultWorkspaceId, token.Token, toggl.SubscriptionCreateData{
		Enabled:     false,
		UrlCallback: url,
		Description: "Какое-то описание",
		EventFilters: []toggl.SubscriptionEventData{
			{
				Action: "*",
				Entity: "time_entry",
			},
		},
	})

	if err != nil {
		panic(err)
	}
}
