package gitclient

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
			BranchName:     event.ObjectAttributes.TargetBranch,
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
		subscribes := database.GetSubscribesByProjectIdAndKind(database.GetSubscribesFilter{
			ProjectId:      event.Project.ID,
			Event:          event.ObjectKind,
			Status:         event.ObjectAttributes.Status,
			AuthorUsername: event.User.Username,
			BranchName:     event.ObjectAttributes.Ref,
		})

		var data interface{}

		data = PipelineDefaultType{
			Event: event,
		}

		if event.ObjectAttributes.Status == "success" {
			if event.ObjectAttributes.Ref == "develop" {
				commits, err := GetCommitsLastPipeline(event.Project.ID, event.ObjectAttributes.BeforeSHA, event.ObjectAttributes.SHA)

				if err != nil {
					break
				}

				data = PipelineCommitsType{
					PipelineDefaultType: data.(PipelineDefaultType),
					Commits:             commits,
				}
			} else if event.ObjectAttributes.Ref == "master" || event.ObjectAttributes.Ref == "release" {
				commits, err := GetCommitsLastPipeline(event.Project.ID, event.ObjectAttributes.BeforeSHA, event.ObjectAttributes.SHA)

				if err != nil {
					break
				}

				data = PipelineLogType{
					PipelineDefaultType: data.(PipelineDefaultType),
					Commits:             commits,
				}
			}
		}

		for _, subscribe := range subscribes {
			var keyboard *tgbotapi.InlineKeyboardMarkup
			switch data := data.(type) {
			case PipelineDefaultType:
				data.Subscribe = &subscribe
				message = data.Make(false)
				keyboard = data.Keyboard()
				break
			case PipelineCommitsType:
				data.Subscribe = &subscribe
				message = data.Make(false)
				keyboard = data.Keyboard()
				break
			case PipelineLogType:
				data.Subscribe = &subscribe
				message = data.Make(false)
				keyboard = data.Keyboard()
				break
			default:
				message = ""
				keyboard = nil
			}

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
						switch data := data.(type) {
						case PipelineCommitsType:
							data.WithPipelineButton = true
							message = data.Make(true)
							keyboard = data.Keyboard()
							telegram.SendMessage(&subscribe.TelegramChannel, message, keyboard, nil)
							break
						case PipelineLogType:
							data.WithPipelineButton = true
							message = data.Make(true)
							keyboard = data.Keyboard()
							telegram.SendMessage(&subscribe.TelegramChannel, message, keyboard, nil)
							break
						}
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
