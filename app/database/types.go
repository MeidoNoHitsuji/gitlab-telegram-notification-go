package database

type GetSubscribesFilter struct {
	ProjectId      int    `json:"project_id"`
	Event          string `json:"event"`
	Status         string `json:"status"`
	AuthorUsername string `json:"author_username"`
	ToBranchName   string `json:"to_branch_name"`
	FromBranchName string `json:"from_branch_name"`
	IsMerge        string `json:"is_merge"`
}
