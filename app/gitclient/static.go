package gitclient

import (
	"fmt"
	"github.com/xanzy/go-gitlab"
	"os"
)

func Subscribe(project *gitlab.Project, hookOptions gitlab.AddProjectHookOptions) (string, error) {
	git := Instant()

	//TODO: Пофиксить тут пагинатор
	hooks, _, err := git.Projects.ListProjectHooks(project.ID, &gitlab.ListProjectHooksOptions{
		Page:    1,
		PerPage: 100,
	})

	if err != nil {
		return "", err
	}

	webhookUrl := fmt.Sprintf("%s/%s", os.Getenv("WEBHOOK_DOMAIN"), os.Getenv("WEBHOOK_URL"))

	hook := gitlab.ProjectHook{
		ID: 0,
	}

	for _, h := range hooks {
		if webhookUrl == h.URL {
			hook = *h
		}
	}

	hookOptions.Token = gitlab.String(os.Getenv("GITLAB_TOKEN"))
	hookOptions.URL = gitlab.String(webhookUrl)

	var text string

	if hook.ID == 0 {
		_, _, err := git.Projects.AddProjectHook(project.ID, &hookOptions)

		if err != nil {
			return "", err
		}

		text = fmt.Sprintf("Подписка на проект [%s](%s) (%d) была добавлена.", project.Name, project.WebURL, project.ID)
	} else {
		_, _, err := git.Projects.EditProjectHook(project.ID, hook.ID, &gitlab.EditProjectHookOptions{
			ConfidentialIssuesEvents: hookOptions.ConfidentialIssuesEvents,
			ConfidentialNoteEvents:   hookOptions.ConfidentialNoteEvents,
			DeploymentEvents:         hookOptions.DeploymentEvents,
			EnableSSLVerification:    hookOptions.EnableSSLVerification,
			IssuesEvents:             hookOptions.IssuesEvents,
			JobEvents:                hookOptions.JobEvents,
			MergeRequestsEvents:      hookOptions.MergeRequestsEvents,
			NoteEvents:               hookOptions.NoteEvents,
			PipelineEvents:           hookOptions.PipelineEvents,
			PushEvents:               hookOptions.PushEvents,
			PushEventsBranchFilter:   hookOptions.PushEventsBranchFilter,
			ReleasesEvents:           hookOptions.ReleasesEvents,
			TagPushEvents:            hookOptions.TagPushEvents,
			Token:                    hookOptions.Token,
			WikiPageEvents:           hookOptions.WikiPageEvents,
			URL:                      hookOptions.URL,
		})

		if err != nil {
			return "", err
		}

		text = fmt.Sprintf("Подписка на проект [%s](%s) (%d) была обновлена.", project.Name, project.WebURL, project.ID)
	}

	return text, nil
}

func DropSubscribe(project *gitlab.Project) {

}
