package models

type TelegramChannel struct {
	ID         int64       `json:"id"`
	Subscribes []Subscribe `json:"subscribes"`
}
