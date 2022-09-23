package toggl

import "time"

type UserType struct {
	Id                 int         `json:"id"`
	ApiToken           string      `json:"api_token"`
	Email              string      `json:"email"`
	FullName           string      `json:"fullname"`
	Timezone           string      `json:"timezone"`
	DefaultWorkspaceId int         `json:"default_workspace_id"`
	BeginningOfWeek    int         `json:"beginning_of_week"`
	ImageUrl           string      `json:"image_url"`
	CreatedAt          time.Time   `json:"created_at"`
	UpdatedAt          time.Time   `json:"updated_at"`
	OpenidEmail        interface{} `json:"openid_email"`
	OpenidEnabled      bool        `json:"openid_enabled"`
	CountryId          int         `json:"country_id"`
	At                 time.Time   `json:"at"`
	IntercomHash       string      `json:"intercom_hash"`
	OauthProviders     []string    `json:"oauth_providers"`
	HasPassword        bool        `json:"has_password"`
}

type Event struct {
	Action   string
	Entities []string
}

type SubscriptionEventData struct {
	Action string `json:"action"`
	Entity string `json:"entity"`
}

type SubscriptionCreateData struct {
	CreatedAt        time.Time               `json:"created_at"`
	DeletedAt        time.Time               `json:"deleted_at"`
	Description      string                  `json:"description"`
	Enabled          bool                    `json:"enabled"`
	EventFilters     []SubscriptionEventData `json:"event_filters"`
	HasPendingEvents bool                    `json:"has_pending_events"`
	Secret           string                  `json:"secret"`
	SubscriptionId   int                     `json:"subscription_id"`
	UpdatedAt        time.Time               `json:"updated_at"`
	UrlCallback      string                  `json:"url_callback"`
	UserId           int                     `json:"user_id"`
	ValidatedAt      time.Time               `json:"validated_at"`
	WorkspaceId      int                     `json:"workspace_id"`
}
