package models

type ToggleJiraIntegration struct {
	ID uint `gorm:"primarykey"`

	TimeEntityId    int64 `json:"time_entity_id"`
	IssueId         int   `json:"issue_id"`
	WorklogRecordId int   `json:"worklog_record_id"`
}
