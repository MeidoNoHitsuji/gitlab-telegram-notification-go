package gitclient

import (
	"fmt"
	"github.com/xanzy/go-gitlab"
	"gitlab-telegram-notification-go/database"
	"gitlab-telegram-notification-go/telegram"
	"os"
	"strings"
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

	hookOptions.Token = gitlab.String(os.Getenv("GITLAB_SECRET"))
	hookOptions.URL = gitlab.String(webhookUrl)

	var text string

	if hook.ID == 0 {
		if os.Getenv("WEBHOOK_TEST") != "true" {
			_, _, err := git.Projects.AddProjectHook(project.ID, &hookOptions)

			if err != nil {
				return "", err
			}
		}

		text = fmt.Sprintf("📝 | Подписка на проект [%s](%s) (%d) была добавлена.", project.Name, project.WebURL, project.ID)
	} else {
		if os.Getenv("WEBHOOK_TEST") != "true" {
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
		}

		text = fmt.Sprintf("📝 | Подписка на проект [%s](%s) (%d) была обновлена.", project.Name, project.WebURL, project.ID)
	}

	return text, nil
}

func Handler(event interface{}) error {
	switch event := event.(type) {
	case *gitlab.MergeEvent:
		subscribes := database.GetSubscribesByProjectIdAndKind(event.Project.ID, event.ObjectKind)
		var message string

		if event.ObjectAttributes.MergeStatus == "unchecked" {
			message = fmt.Sprintf("🎭⚠ Необходимо проверить MergeRequest! | [%s](%s) (%d)", event.Project.Name, event.Project.WebURL, event.Project.ID)
		} else if event.ObjectAttributes.MergeStatus == "cannot_be_merged" {
			message = fmt.Sprintf("🎭❌ Обнаружены ошибки в MergeRequest! | [%s](%s) (%d)", event.Project.Name, event.Project.WebURL, event.Project.ID)
		} else if event.ObjectAttributes.MergeStatus == "can_be_merged" {
			message = fmt.Sprintf("🎭✅ Был завершён MergeRequest! | [%s](%s) (%d)", event.Project.Name, event.Project.WebURL, event.Project.ID)
		} else {
			break
		}

		message = fmt.Sprintf("%s\n—————\n[%s](%s)", message, event.ObjectAttributes.Title, event.ObjectAttributes.URL)
		message = fmt.Sprintf("%s\n\n🌳: %s → %s", message, event.ObjectAttributes.SourceBranch, event.ObjectAttributes.TargetBranch)
		message = fmt.Sprintf("%s\n🧙: [%s](%s/%s)", message, event.User.Name, os.Getenv("GITLAB_URL"), event.User.Username)

		for _, subscribe := range subscribes {
			telegram.SendMessage(&subscribe.TelegramChannel, message, nil)
		}
	case *gitlab.PipelineEvent:
		subscribes := database.GetSubscribesByProjectIdAndKind(event.Project.ID, event.ObjectKind)
		var message string
		if event.ObjectAttributes.Status == "failed" {
			message = fmt.Sprintf("🧩❌ PipeLine завершился ошибкой! | [%s](%s) (%d)", event.Project.Name, event.Project.WebURL, event.Project.ID)
			message = fmt.Sprintf("%s\n—————", message)
		} else if event.ObjectAttributes.Status == "success" {
			message = fmt.Sprintf("🧩✅ PipeLine завершился успешно! | [%s](%s) (%d)", event.Project.Name, event.Project.WebURL, event.Project.ID)
			message = fmt.Sprintf("%s\n—————", message)
		} else {
			break
		}

		if event.MergeRequest.ID != 0 {
			message = fmt.Sprintf("%s\n[%s](%s/-/pipelines/%d)\n—————", message, event.MergeRequest.Title, event.Project.WebURL, event.ObjectAttributes.ID)
		} else {
			messages := strings.Split(event.Commit.Message, "\n")

			if len(messages) > 0 {
				message = fmt.Sprintf("%s\n[%s](%s/-/pipelines/%d)\n—————", message, messages[0], event.Project.WebURL, event.ObjectAttributes.ID)
			} else {
				message = fmt.Sprintf("%s\n[%s](%s/-/pipelines/%d)\n—————", message, event.Commit.Message, event.Project.WebURL, event.ObjectAttributes.ID)
			}
		}

		message = fmt.Sprintf("%s\nСборочная линия:", message)

		for _, stage := range event.ObjectAttributes.Stages {
			for _, build := range event.Builds {
				if build.Stage == stage {
					emoji := "❓"

					if build.Status == "failed" {
						emoji = "❌"
					} else if build.Status == "skipped" {
						emoji = "⏩"
					} else if build.Status == "success" {
						emoji = "✅"
					}

					message = fmt.Sprintf("%s\n%s [%s](%s/-/jobs/%d)", message, emoji, build.Name, event.Project.WebURL, build.ID)
				}
			}
		}

		message = fmt.Sprintf("%s\n", message)

		if event.MergeRequest.ID != 0 {
			message = fmt.Sprintf("%s\n🌳: %s → %s", message, event.MergeRequest.SourceBranch, event.MergeRequest.TargetBranch)
		} else {
			message = fmt.Sprintf("%s\n🌳: %s", message, event.ObjectAttributes.Ref)
		}

		message = fmt.Sprintf("%s\n🧙: [%s](%s/%s)", message, event.User.Name, os.Getenv("GITLAB_URL"), event.User.Username)

		for _, subscribe := range subscribes {
			telegram.SendMessage(&subscribe.TelegramChannel, message, nil)
		}
	}
	return nil
}
