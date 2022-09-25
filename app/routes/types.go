package routes

import "time"

type ToggleData struct {
	EventId   int       `json:"event_id"`
	CreatedAt time.Time `json:"created_at"`
	CreatorId int       `json:"creator_id"`
	Metadata  struct {
		RequestType string `json:"request_type"`
		EventUserId int    `json:"event_user_id"`
	} `json:"metadata"`
	Payload           string `json:"payload"`
	SubscriptionId    int    `json:"subscription_id"`
	ValidationCode    string `json:"validation_code"`
	ValidationCodeUrl string `json:"validation_code_url"`
}
