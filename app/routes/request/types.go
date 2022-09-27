package request

import "time"

type ValidationData struct {
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

type ToggleData struct {
	EventId   int64     `json:"event_id"`
	CreatedAt time.Time `json:"created_at"`
	CreatorId int       `json:"creator_id"`
	Metadata  struct {
		Action           string `json:"action"`
		EventUserId      string `json:"event_user_id"`
		Model            string `json:"model"`
		ModelOwnerId     string `json:"model_owner_id"`
		Path             string `json:"path"`
		ProjectId        string `json:"project_id"`
		ProjectIsPrivate string `json:"project_is_private"`
		RequestType      string `json:"request_type"`
		TimeEntryId      string `json:"time_entry_id"`
		WorkspaceId      string `json:"workspace_id"`
	} `json:"metadata"`
	Payload struct {
		At              time.Time     `json:"at"`
		Billable        bool          `json:"billable"`
		Description     string        `json:"description"`
		Duration        int           `json:"duration"`
		Duronly         bool          `json:"duronly"`
		Id              int64         `json:"id"`
		Pid             int           `json:"pid"`
		ProjectId       int           `json:"project_id"`
		ServerDeletedAt interface{}   `json:"server_deleted_at"`
		Start           time.Time     `json:"start"`
		Stop            interface{}   `json:"stop"`
		TagIds          interface{}   `json:"tag_ids"`
		Tags            []interface{} `json:"tags"`
		TaskId          interface{}   `json:"task_id"`
		Uid             int           `json:"uid"`
		UserId          int           `json:"user_id"`
		Wid             int           `json:"wid"`
		WorkspaceId     int           `json:"workspace_id"`
	} `json:"payload"`
	SubscriptionId int `json:"subscription_id"`
}
