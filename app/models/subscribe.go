package models

import "gorm.io/gorm"

type Subscribe struct {
	gorm.Model
	ProjectId         uint             `json:"not null;project_id"`
	TelegramChannelId int64            `json:"not null;telegram_channel_id"`
	Project           Project          `gorm:"foreignKey:ProjectId;references:ID;" json:"project"`
	TelegramChannel   TelegramChannel  `gorm:"foreignKey:TelegramChannelId;references:ID;" json:"telegram_channel"`
	Events            []SubscribeEvent `json:"events"`
}

type SubscribeEvent struct {
	SubscribeId uint      `json:"not null;subscribe_id"`
	Event       string    `json:"not null;event"`
	Subscribe   Subscribe `gorm:"foreignKey:SubscribeId;references:ID;" json:"subscribe"`
}
