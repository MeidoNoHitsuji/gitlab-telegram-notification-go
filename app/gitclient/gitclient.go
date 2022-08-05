package gitclient

import (
	"github.com/xanzy/go-gitlab"
	"log"
	"os"
)

var instant *gitlab.Client

func New(GitlabToken string, GitlabUrl string) *gitlab.Client {
	git, err := gitlab.NewClient(GitlabToken, gitlab.WithBaseURL(GitlabUrl))

	if err != nil {
		log.Fatalf("Failed to create gitlab telegram: %v", err)
	}

	return git
}

func Instant() *gitlab.Client {
	if instant == nil {
		instant = New(os.Getenv("GITLAB_TOKEN"), os.Getenv("GITLAB_URL"))
	}

	return instant
}
