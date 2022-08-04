package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username          string          `json:"username"`
	TelegramChannelId int64           `json:"telegram_channel_id"`
	TelegramChannel   TelegramChannel `gorm:"foreignKey:TelegramChannelId;references:ID;" json:"telegram_channel"`
}
