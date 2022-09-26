package cmd

import (
	"github.com/spf13/cobra"
	"gitlab-telegram-notification-go/database"
	"gitlab-telegram-notification-go/models"
)

var seedCmd = &cobra.Command{
	Use:  "seed",
	RunE: seed,
}

func init() {
	rootCmd.AddCommand(seedCmd)
}

func seed(cmd *cobra.Command, args []string) error {
	db := database.Instant()

	var events []models.SubscribeEvent

	db.Model(models.SubscribeEvent{}).Where("parameters is null").Preload("Subscribe").Find(&events)

	for _, event := range events {

		if event.Event == "pipeline" {
			event.Parameters = map[string][]string{
				"to_branch_name": {
					"develop",
					"release",
					"master",
				},
				"is_merge": {"false"},
				"status":   {"success"},
			}
			event.Formatter = "logs"
			db.Omit("Subscribe").Save(&event)
			newEvent := models.SubscribeEvent{
				SubscribeId: event.SubscribeId,
				Event:       event.Event,
				Parameters: map[string][]string{
					"is_merge": {"false"},
					"status":   {"failed"},
				},
				Formatter: "default",
			}
			db.Omit("Subscribe").Create(&newEvent)
		} else {
			event.Parameters = map[string][]string{}
			db.Omit("Subscribe").Save(&event)
		}
	}

	return nil
}
