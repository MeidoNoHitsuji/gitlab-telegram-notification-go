package gitclient

import (
	"encoding/json"
	"fmt"
	"github.com/xanzy/go-gitlab"
	"gitlab-telegram-notification-go/database"
	"gitlab-telegram-notification-go/helper"
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

func SubscribeByProject(project *gitlab.Project) (string, error) {
	allEvents := database.GetEventsByProjectId(project.ID)

	return Subscribe(project, gitlab.AddProjectHookOptions{
		PushEvents:          gitlab.Bool(helper.Contains(allEvents, helper.Slugify(string(gitlab.EventTypePush)))),
		PipelineEvents:      gitlab.Bool(helper.Contains(allEvents, helper.Slugify(string(gitlab.EventTypePipeline)))),
		MergeRequestsEvents: gitlab.Bool(helper.Contains(allEvents, helper.Slugify(string(gitlab.EventTypeMergeRequest)))),
	})
}

func Handler(event interface{}) error {
	switch event := event.(type) {
	case *gitlab.MergeEvent:
		subscribeEvents := database.GetSubscribesByProjectIdAndKind(database.GetSubscribesFilter{
			ProjectId:      event.Project.ID,
			Event:          event.ObjectKind,
			Status:         event.ObjectAttributes.State,
			AuthorUsername: event.User.Username,
			ToBranchName:   event.ObjectAttributes.TargetBranch,
		})
		var message string

		for _, subscribeEvent := range subscribeEvents {

			data := MergeDefaultType{
				Event:     event,
				Subscribe: &subscribeEvent.Subscribe,
			}

			message = data.Make()

			if message == "" {
				continue
			}

			telegram.SendMessageById(subscribeEvent.Subscribe.TelegramChannelId, message, nil, nil)
		}
	case *gitlab.PipelineEvent:
		var message string
		var IsMerge string

		if event.MergeRequest.ID != 0 {
			IsMerge = "true"
		} else {
			IsMerge = "false"
		}

		subscribeEvents := database.GetSubscribesByProjectIdAndKind(database.GetSubscribesFilter{
			ProjectId:      event.Project.ID,
			Event:          event.ObjectKind,
			Status:         event.ObjectAttributes.Status,
			AuthorUsername: event.User.Username,
			ToBranchName:   event.ObjectAttributes.Ref,
			IsMerge:        IsMerge,
		})

		var data PipelineDefaultInterface

		data = NewPipelineDefaultType(event)
		beforePipeline, err := GetBeforeFailedPipeline(event.Project.ID, event.ObjectAttributes.BeforeSHA, event.ObjectAttributes.Ref)
		beforeSHA := event.ObjectAttributes.BeforeSHA

		if err == nil && beforePipeline != nil {
			beforeSHA = beforePipeline.BeforeSHA
		}

		commits, err := GetCommitsLastPipeline(event.Project.ID, beforeSHA, event.ObjectAttributes.SHA)

		if err != nil {
			break
		}

		for _, subscribeEvent := range subscribeEvents {

			if subscribeEvent.Formatter == "commits" {
				data = NewPipelineCommitsType(event, commits)
			} else if subscribeEvent.Formatter == "logs" {
				data = NewPipelineLogType(event, commits)
			}

			data.SetSubscribe(&subscribeEvent.Subscribe)

			message = data.Make(false)
			keyboard := data.Keyboard(false)

			if message == "" {
				continue
			}

			_, err := telegram.SendMessageById(subscribeEvent.Subscribe.TelegramChannelId, message, keyboard, nil)

			var errMap map[string]interface{}
			out, _ := json.Marshal(err)

			err = json.Unmarshal(out, &errMap)

			if err == nil {
				code, ok := errMap["Code"]
				if ok {
					if code.(float64) == 400 {
						message = data.Make(true)
						keyboard = data.Keyboard(true)
						telegram.SendMessageById(subscribeEvent.Subscribe.TelegramChannelId, message, keyboard, nil)
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

func GetBeforeFailedPipeline(projectId int, sha string, ref string) (*gitlab.Pipeline, error) {
	git := Instant()
	result, _, err := git.Pipelines.ListProjectPipelines(340, &gitlab.ListProjectPipelinesOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 1,
		},
		Status: gitlab.BuildState(gitlab.Failed),
		SHA:    gitlab.String(sha),
		Ref:    gitlab.String(ref),
	})

	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, nil
	} else {
		pipeline, _, err := git.Pipelines.GetPipeline(projectId, result[0].ID)

		if err != nil {
			return nil, err
		}

		beforePipeline, err := GetBeforeFailedPipeline(projectId, pipeline.BeforeSHA, pipeline.Ref)

		if err == nil && beforePipeline != nil {
			return beforePipeline, err
		} else {
			return pipeline, nil
		}
	}
}
