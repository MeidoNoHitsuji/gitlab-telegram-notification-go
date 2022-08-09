package gitclient

import (
	"fmt"
	"github.com/xanzy/go-gitlab"
	"gitlab-telegram-notification-go/database"
	"gitlab-telegram-notification-go/telegram"
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

	webhookUrl := fmt.Sprintf("%s:%s/%s", os.Getenv("WEBHOOK_DOMAIN"), os.Getenv("WEBHOOK_PORT"), os.Getenv("WEBHOOK_URL"))

	hook := gitlab.ProjectHook{
		ID: 0,
	}

	for _, h := range hooks {
		if webhookUrl == h.URL {
			hook = *h
		}
	}

	hookOptions.Token = gitlab.String(os.Getenv("GITLAB_SECRET"))
	hookOptions.URL = gitlab.String(webhookUrl)

	var text string

	if hook.ID == 0 {
		_, _, err := git.Projects.AddProjectHook(project.ID, &hookOptions)

		if err != nil {
			return "", err
		}

		text = fmt.Sprintf("📝 | Подписка на проект [%s](%s) (%d) была добавлена.", project.Name, project.WebURL, project.ID)
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

		text = fmt.Sprintf("📝 | Подписка на проект [%s](%s) (%d) была обновлена.", project.Name, project.WebURL, project.ID)
	}

	return text, nil
}

func Handler(event interface{}) error {
	switch event := event.(type) {
	case *gitlab.MergeEvent:
		subscribes := database.GetSubscribesByProjectIdAndKind(event.Project.ID, event.ObjectKind)

		for _, subscribe := range subscribes {
			message := fmt.Sprintf("🎭 Создан новый MergeRequest! | [%s](%s) %d", event.Project.Name, event.Project.WebURL, event.Project.ID)
			message = fmt.Sprintf("%s\n----------\n[%s](%s)", message, event.ObjectAttributes.Title, event.ObjectAttributes.URL)
			message = fmt.Sprintf("%s\n\n🌳: %s -> %s", message, event.ObjectAttributes.SourceBranch, event.ObjectAttributes.TargetBranch)
			message = fmt.Sprintf("%s\n🧙: [%s](%s/%s)", message, event.User.Name, os.Getenv("GITLAB_URL"), event.User.Username)
			telegram.SendMessage(&subscribe.TelegramChannel, message)
		}
	case *gitlab.PipelineEvent:
		subscribes := database.GetSubscribesByProjectIdAndKind(event.Project.ID, event.ObjectKind)

		for _, subscribe := range subscribes {
			telegram.SendMessage(&subscribe.TelegramChannel, "Закончился новый PipeLine!")
		}
	}
	return nil
}
