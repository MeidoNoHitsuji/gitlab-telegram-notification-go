package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab-telegram-notification-go/database"
	"gitlab-telegram-notification-go/models"
	"gitlab-telegram-notification-go/toggl"
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
		TokenType: models.ToggleToken,
		User: models.User{
			TelegramChannelId: 479413765,
		},
	}).First(&token)

	userData, err := toggl.Me(token.Token)

	if err != nil {
		panic(err)
	}

	r, err := toggl.GetSubscriptions(userData.DefaultWorkspaceId, token.Token)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(r)
	}

	for _, data := range r {
		res, err := toggl.EnableSubscriptions(userData.DefaultWorkspaceId, data.SubscriptionId, true, token.Token)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(res)
		}
	}

}
