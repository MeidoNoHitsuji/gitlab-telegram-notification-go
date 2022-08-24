package database

type GetSubscribesFilter struct {
	ProjectId      int    `json:"project_id"`
	Event          string `json:"event"`
	Status         string `json:"status"`
	AuthorUsername string `json:"author_username"`
	BranchName     string `json:"branch_name"`
}
