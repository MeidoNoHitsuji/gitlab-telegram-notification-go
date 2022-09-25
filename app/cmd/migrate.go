package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab-telegram-notification-go/database"
	"gitlab-telegram-notification-go/models"
)

var migrateCmd = &cobra.Command{
	Use:  "migrate",
	RunE: migrate,
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}

func migrate(cmd *cobra.Command, args []string) error {
	db := database.Instant()

	db.Exec("ALTER TABLE `subscribe_events` ADD `id` bigint unsigned primary key AUTO_INCREMENT")

	ms := []interface{}{
		&models.TelegramChannel{},
		&models.User{},
		&models.Project{},
		&models.Subscribe{},
		&models.SubscribeEvent{},
		&models.UserTelegramChannelAction{},
		&models.UserToken{},
		models.UserIntegrations{},
	}

	if err := db.AutoMigrate(ms...); err != nil {
		return err
	}

	fmt.Println("Migrate completed.")

	return nil
}
