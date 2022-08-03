package client

import (
	"github.com/xanzy/go-gitlab"
	"log"
	"os"
)

var gitlabInstant *gitlab.Client

func Gitlab() *gitlab.Client {
	if gitlabInstant != nil {
		return gitlabInstant
	}

	git, err := gitlab.NewClient(os.Getenv("GITLAB_TOKEN"), gitlab.WithBaseURL(os.Getenv("GITLAB_URL")))

	if err != nil {
		log.Fatalf("Failed to create gitlab client: %v", err)
	}

	gitlabInstant = git
	return gitlabInstant
}
