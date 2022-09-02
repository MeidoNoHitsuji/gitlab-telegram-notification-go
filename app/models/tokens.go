package models

import "gorm.io/gorm"

const (
	ToggleTrackerTokenType = "toggle_tracker"
	JiraTokenType          = "jira"
)

type UserToken struct {
	gorm.Model
	Token     string `json:"token"`
	TokenType string `json:"token_type"`
	UserId    uint   `json:"user_id"`
	User      User   `gorm:"foreignKey:UserId;references:ID;" json:"user"`
}
