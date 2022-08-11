package models

import (
	"gorm.io/gorm"
)

type Subscribe struct {
	gorm.Model
	ProjectId         int              `json:"not null;project_id"`
	TelegramChannelId int64            `json:"not null;telegram_channel_id"`
	Project           Project          `gorm:"foreignKey:ProjectId;references:ID;" json:"project"`
	TelegramChannel   TelegramChannel  `gorm:"foreignKey:TelegramChannelId;references:ID;" json:"telegram_channel"`
	Active            bool             `gorm:"default:true" json:"active"`
	Events            []SubscribeEvent `json:"events"`
}

type SubscribeEvent struct {
	SubscribeId uint      `json:"subscribe_id"`
	Event       string    `json:"event"`
	Subscribe   Subscribe `gorm:"foreignKey:SubscribeId;references:ID;" json:"subscribe"`
}
