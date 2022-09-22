package models

import (
	"encoding/json"
	"gorm.io/gorm"
)

type JSON json.RawMessage

type Subscribe struct {
	gorm.Model
	ProjectId         int              `json:"project_id"`
	TelegramChannelId int64            `json:"telegram_channel_id"`
	Project           Project          `gorm:"foreignKey:ProjectId;references:ID;" json:"project"`
	TelegramChannel   TelegramChannel  `gorm:"foreignKey:TelegramChannelId;references:ID;" json:"telegram_channel"`
	Events            []SubscribeEvent `json:"events"`
}

type SubscribeEvent struct {
	ID          uint                `gorm:"primarykey"`
	SubscribeId uint                `json:"subscribe_id"`
	Event       string              `json:"event"`
	Parameters  map[string][]string `gorm:"serializer:json"`
	Formatter   string              `json:"formatter"`
	Subscribe   Subscribe           `gorm:"foreignKey:SubscribeId;references:ID;" json:"subscribe"`
}
