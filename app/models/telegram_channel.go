package models

type TelegramChannel struct {
	ID         int64       `json:"id"`
	Active     bool        `json:"active"`
	Subscribes []Subscribe `json:"subscribes"`
}
