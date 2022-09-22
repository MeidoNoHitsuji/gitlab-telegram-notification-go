package gitclient

import (
	"encoding/json"
	"fmt"
	"github.com/xanzy/go-gitlab"
	"gitlab-telegram-notification-go/database"
	fm "gitlab-telegram-notification-go/helper/formater"
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

	port := os.Getenv("WEBHOOK_PORT")

	var webhookUrl string

	if port != "" {
		webhookUrl = fmt.Sprintf("%s:%s/%s", os.Getenv("WEBHOOK_DOMAIN"), port, os.Getenv("WEBHOOK_URL"))
	} else {
		webhookUrl = fmt.Sprintf("%s/%s", os.Getenv("WEBHOOK_DOMAIN"), os.Getenv("WEBHOOK_URL"))
	}

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

		text = fmt.Sprintf("📝 | Подписка на проект %s (%d) была добавлена.", fm.Link(project.Name, project.WebURL), project.ID)
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

		text = fmt.Sprintf("📝 | Подписка на проект %s (%d) была обновлена.", fm.Link(project.Name, project.WebURL), project.ID)
	}

	return text, nil
}

func Handler(event interface{}) error {
	switch event := event.(type) {
	case *gitlab.MergeEvent:
		subscribes := database.GetSubscribesByProjectIdAndKind(database.GetSubscribesFilter{
			ProjectId:      event.Project.ID,
			Event:          event.ObjectKind,
			Status:         event.ObjectAttributes.State,
			AuthorUsername: event.User.Username,
			ToBranchName:   event.ObjectAttributes.TargetBranch,
		})
		var message string

		for _, subscribe := range subscribes {

			data := MergeDefaultType{
				Event:     event,
				Subscribe: &subscribe,
			}

			message = data.Make()

			if message == "" {
				continue
			}

			telegram.SendMessage(&subscribe.TelegramChannel, message, nil, nil)
		}
	case *gitlab.PipelineEvent:
		var message string
		var IsMerge string

		if event.MergeRequest.ID != 0 {
			IsMerge = "true"
		} else {
			IsMerge = "false"
		}

		subscribes := database.GetSubscribesByProjectIdAndKind(database.GetSubscribesFilter{
			ProjectId:      event.Project.ID,
			Event:          event.ObjectKind,
			Status:         event.ObjectAttributes.Status,
			AuthorUsername: event.User.Username,
			ToBranchName:   event.ObjectAttributes.Ref,
			IsMerge:        IsMerge,
		})

		var data PipelineDefaultInterface

		data = NewPipelineDefaultType(event)
		commits, err := GetCommitsLastPipeline(event.Project.ID, event.ObjectAttributes.BeforeSHA, event.ObjectAttributes.SHA)

		if err != nil {
			break
		}

		for _, subscribe := range subscribes {

			if len(subscribe.Events) == 0 {
				continue
			}

			subscribeEvent := subscribe.Events[0]

			if subscribeEvent.Formatter == "commit" {
				data = NewPipelineCommitsType(event, commits)
			} else if subscribeEvent.Formatter == "log" {
				data = NewPipelineLogType(event, commits)
			}

			data.SetSubscribe(&subscribe)

			message = data.Make(false)
			keyboard := data.Keyboard(false)

			if message == "" {
				continue
			}

			_, err := telegram.SendMessage(&subscribe.TelegramChannel, message, keyboard, nil)

			var errMap map[string]interface{}
			out, _ := json.Marshal(err)

			err = json.Unmarshal(out, &errMap)

			if err == nil {
				code, ok := errMap["Code"]
				if ok {
					if code.(float64) == 400 {
						message = data.Make(true)
						keyboard = data.Keyboard(true)
						telegram.SendMessage(&subscribe.TelegramChannel, message, keyboard, nil)
						break
					}
				}
			} else {
				fmt.Println(err)
			}
		}
	}
	return nil
}

func GetCommitsLastPipeline(projectId int, fromHash string, toHash string) ([]*gitlab.Commit, error) {
	git := Instant()

	compare, _, err := git.Repositories.Compare(projectId, &gitlab.CompareOptions{
		From: gitlab.String(fromHash),
		To:   gitlab.String(toHash),
	})

	if err != nil {
		return nil, err
	}

	return compare.Commits, nil
}
