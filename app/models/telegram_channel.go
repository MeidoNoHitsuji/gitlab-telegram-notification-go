package models

type TelegramChannel struct {
	ID         int64       `json:"id"`
	Active     bool        `json:"active"`
	Subscribes []Subscribe `json:"subscribes"`
}

type UserTelegramChannelAction struct {
	UserId            uint            `json:"user_id"`
	TelegramChannelId int64           `json:"telegram_channel_id"`
	Action            string          `json:"action"`
	User              User            `gorm:"foreignKey:UserId;references:ID;" json:"user"`
	TelegramChannel   TelegramChannel `gorm:"foreignKey:TelegramChannelId;references:ID;" json:"telegram_channel"`
}
