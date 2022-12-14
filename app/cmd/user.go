package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab-telegram-notification-go/database"
	"gitlab-telegram-notification-go/models"
)

var userCmd = &cobra.Command{
	Use:  "user",
	Run:  user,
	Args: cobra.MinimumNArgs(1),
}

func init() {
	rootCmd.AddCommand(userCmd)
}

func user(cmd *cobra.Command, args []string) {
	userId := args[0]
	var user models.User

	db := database.Instant()
	result := db.Model(&models.User{}).Where("telegram_channel_id = ?", userId).Find(&user)

	if result.RowsAffected == 0 {
		fmt.Println("Users not found!")
	} else {
		fmt.Printf("Username: %s", user.Username)
	}
}
