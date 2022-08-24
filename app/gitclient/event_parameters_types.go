package gitclient

type PipelineParameters struct {
	AuthorUsername string `json:"author_username"`
	BranchName     string `json:"branch_name"`
	Status         string `json:"status"`
}
