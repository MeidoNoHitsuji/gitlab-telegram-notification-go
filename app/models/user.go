package models

import "gorm.io/gorm"

const (
	ToggleJiraIntegration = "toggle_jira_integration"
)

type User struct {
	gorm.Model
	Username          string          `json:"username"`
	TelegramChannelId int64           `json:"telegram_channel_id"`
	TelegramChannel   TelegramChannel `gorm:"foreignKey:TelegramChannelId;references:ID;" json:"telegram_channel"`
}

type UserIntegrations struct {
	ID              uint   `gorm:"primarykey"`
	IntegrationType string `json:"integration_type"`
	Active          bool   `json:"active"`
	UserId          uint   `json:"user_id"`
	User            User   `gorm:"foreignKey:UserId;references:ID;" json:"user"`
}
