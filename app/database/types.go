package database

type GetSubscribesFilter struct {
	ProjectId         int    `json:"project_id"`
	Event             string `json:"event"`
	Status            string `json:"status"`
	AuthorUsername    string `json:"author_username"`
	NotAuthorUsername string `json:"not_author_username"`
	Source            string `json:"source"`
	ToBranchName      string `json:"to_branch_name"`
	FromBranchName    string `json:"from_branch_name"`
	IsMerge           string `json:"is_merge"`
	State             string `json:"state"`
	Action            string `json:"action"`
}
