package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
	"gitlab-telegram-notification-go/database"
	"gitlab-telegram-notification-go/gitclient"
	"gitlab-telegram-notification-go/helper"
	"gitlab-telegram-notification-go/models"
	"os"
)

var changeWebhookUrlCmd = &cobra.Command{
	Use:  "change-webhook-url",
	Args: cobra.MinimumNArgs(1),
	Run:  changeWebhookUrl,
}

func init() {
	rootCmd.AddCommand(changeWebhookUrlCmd)
}

func changeWebhookUrl(cmd *cobra.Command, args []string) {
	git := gitclient.Instant()
	db := database.Instant()

	newWebhookUrl := args[0]

	var webhookUrl string

	if len(args) > 1 {
		webhookUrl = args[1]
	} else {
		port := os.Getenv("WEBHOOK_PORT")

		if port != "" {
			webhookUrl = fmt.Sprintf("%s:%s/%s", os.Getenv("WEBHOOK_DOMAIN"), port, os.Getenv("WEBHOOK_URL"))
		} else {
			webhookUrl = fmt.Sprintf("%s/%s", os.Getenv("WEBHOOK_DOMAIN"), os.Getenv("WEBHOOK_URL"))
		}
	}

	var projects []models.Project
	db.Find(&projects)

	for _, project := range projects {

		allEvents := database.GetEventsByProjectId(project.ID)

		hookOptions := gitlab.AddProjectHookOptions{
			PushEvents:          gitlab.Bool(helper.Contains(allEvents, helper.Slugify(string(gitlab.EventTypePush)))),
			PipelineEvents:      gitlab.Bool(helper.Contains(allEvents, helper.Slugify(string(gitlab.EventTypePipeline)))),
			MergeRequestsEvents: gitlab.Bool(helper.Contains(allEvents, helper.Slugify(string(gitlab.EventTypeMergeRequest)))),
		}

		hooks, _, err := git.Projects.ListProjectHooks(project.ID, &gitlab.ListProjectHooksOptions{
			Page:    1,
			PerPage: 100,
		})

		if err != nil {
			continue
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
		hookOptions.URL = gitlab.String(newWebhookUrl)

		editProjectHookOptions := gitlab.EditProjectHookOptions{
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
		}

		if os.Getenv("WEBHOOK_TEST") != "true" {
			git.Projects.EditProjectHook(project.ID, hook.ID, &editProjectHookOptions)
		} else {
			out, err := json.Marshal(editProjectHookOptions)

			if err == nil {
				fmt.Println(project.ID, hook.ID, string(out))
				fmt.Println("-----")
			} else {
				fmt.Println(err)
			}
		}
	}
}
