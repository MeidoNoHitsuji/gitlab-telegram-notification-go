package gitclient

type PipelineParameters struct {
	AuthorUsername []string `json:"author_username"`
	ToBranchName   []string `json:"to_branch_name"`
	FromBranchName []string `json:"from_branch_name"`
	Status         []string `json:"status"`
	IsMerge        []string `json:"is_merge"`
}
