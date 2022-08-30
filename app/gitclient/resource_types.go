package gitclient

import (
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/xanzy/go-gitlab"
	"gitlab-telegram-notification-go/gitclient/parser"
	"gitlab-telegram-notification-go/helper"
	fm "gitlab-telegram-notification-go/helper/formater"
	"gitlab-telegram-notification-go/models"
	"os"
	"strings"
)

var types = map[string]string{
	"feat":     "Новое",
	"fix":      "Исправления",
	"chore":    "Рутинные исправления",
	"test":     "Тестирование",
	"build":    "Сборка",
	"refactor": "Рефакторинг кода",
	"docs":     "Обновления документации",
	"ci":       "Изменения CI",
	"perf":     "Исправления производительности",
	"style":    "Декоративные исправления",
	"other":    "Другое",
}

type PipelineDefaultType struct {
	Event     *gitlab.PipelineEvent
	Subscribe *models.Subscribe
}

func (t PipelineDefaultType) Header() (string, error) {
	var message string
	if t.Event.ObjectAttributes.Status == "failed" {
		message = fmt.Sprintf("🧩❌ PipeLine завершился ошибкой! | %s (%d)", fm.Link(t.Event.Project.Name, t.Event.Project.WebURL), t.Event.Project.ID)
		message = fmt.Sprintf("%s\n—————", message)
	} else if t.Event.ObjectAttributes.Status == "success" {
		message = fmt.Sprintf("🧩✅ PipeLine завершился успешно! | %s (%d)", fm.Link(t.Event.Project.Name, t.Event.Project.WebURL), t.Event.Project.ID)
		message = fmt.Sprintf("%s\n—————", message)
	} else {
		return "", errors.New("Такой статус пайплайна не поддерживается.")
	}

	if t.Event.MergeRequest.ID != 0 {
		url := fmt.Sprintf("%s/-/pipelines/%d", t.Event.Project.WebURL, t.Event.ObjectAttributes.ID)
		message = fmt.Sprintf("%s\n%s\n—————", message, fm.Link(t.Event.MergeRequest.Title, url))
	} else {
		messages := strings.Split(t.Event.Commit.Message, "\n")
		url := fmt.Sprintf("%s/-/pipelines/%d", t.Event.Project.WebURL, t.Event.ObjectAttributes.ID)
		if len(messages) > 0 {
			message = fmt.Sprintf("%s\n%s\n—————", message, fm.Link(messages[0], url))
		} else {
			message = fmt.Sprintf("%s\n%s\n—————", message, fm.Link(t.Event.Commit.Message, url))
		}
	}

	return message, nil
}

func (t PipelineDefaultType) Footer() string {
	var message string
	if t.Event.MergeRequest.ID != 0 {
		message = fmt.Sprintf("\n🌳: %s → %s", tgbotapi.EscapeText(tgbotapi.ModeHTML, t.Event.MergeRequest.SourceBranch), tgbotapi.EscapeText(tgbotapi.ModeHTML, t.Event.MergeRequest.TargetBranch))
	} else {
		message = fmt.Sprintf("\n🌳: %s", tgbotapi.EscapeText(tgbotapi.ModeHTML, t.Event.ObjectAttributes.Ref))
	}

	url := fmt.Sprintf("%s/%s", os.Getenv("GITLAB_URL"), t.Event.User.Username)
	message = fmt.Sprintf("%s\n🧙: %s", message, fm.Link(t.Event.User.Name, url))

	return message
}

func (t *PipelineDefaultType) Body() string {
	message := ""

	for _, stage := range t.Event.ObjectAttributes.Stages {
		for _, build := range t.Event.Builds {
			if build.Stage == stage {
				emoji := "❓"

				if build.Status == "failed" {
					emoji = "❌"
				} else if build.Status == "skipped" {
					emoji = "⏩"
				} else if build.Status == "success" {
					emoji = "✅"
				}

				url := fmt.Sprintf("%s/-/jobs/%d", t.Event.Project.WebURL, build.ID)
				message = fmt.Sprintf("%s\n%s %s", message, emoji, fm.Link(build.Name, url))
			}
		}
	}

	return message
}

func (t *PipelineDefaultType) Make() string {
	message, err := t.Header()
	if err != nil {
		return ""
	}
	message = fmt.Sprintf("%s\nСборочная линия:%s", message, t.Body())

	return fmt.Sprintf("%s\n%s", message, t.Footer())
}

type PipelineCommitsType struct {
	PipelineDefaultType
	Commits []*gitlab.Commit
}

func (t *PipelineCommitsType) Body() string {
	message := ""
	for _, commit := range t.Commits {
		if len(commit.ParentIDs) > 1 {
			continue
		}
		message = fmt.Sprintf("%s\n📄 %s", message, fm.Link(commit.Title, commit.WebURL))
	}

	return message
}

func (t *PipelineCommitsType) Make() string {
	message, err := t.Header()
	if err != nil {
		return ""
	}
	message = fmt.Sprintf("%s\nЗалитые изменения:%s", message, t.Body())

	return fmt.Sprintf("%s\n%s", message, t.Footer())
}

type PipelineLogType struct {
	PipelineDefaultType
	Commits []*gitlab.Commit
}

func (t *PipelineLogType) Body() string {
	message := ""
	commits := map[string]map[string][]map[string]interface{}{}

	for _, commit := range t.Commits {
		if len(commit.ParentIDs) > 1 || strings.HasPrefix(commit.Message, "Merge") {
			continue
		}

		resCommit := parser.CompileCommit(commit.Message)

		t := resCommit.Type

		keyTypes := helper.Keys(types)

		if !helper.Contains(keyTypes, t) {
			t = "other"
		}

		_, ok := commits[t]

		if !ok {
			commits[t] = map[string][]map[string]interface{}{}
		}

		scope := resCommit.Scope

		_, ok = commits[t][scope]

		if !ok {
			commits[t][scope] = []map[string]interface{}{}
		}

		body := resCommit.Body

		jira, ok := resCommit.Footer["jira"]

		if !ok {
			jira = []string{}
		}

		commits[t][scope] = append(commits[t][scope], map[string]interface{}{
			"description": resCommit.Description,
			"url":         commit.WebURL,
			"body":        body,
			"jira":        jira,
		})
	}

	for k, v := range types {
		data, ok := commits[k]

		if !ok {
			continue
		}

		subMessage := fmt.Sprintf("\n%s:", fm.Bold(v))
		for scopeKey, dataCommits := range data {
			subMessage = fmt.Sprintf("%s\n    %s:", subMessage, fm.Underline(scopeKey))
			for _, commit := range dataCommits {

				subMessage = fmt.Sprintf("%s\n        📄%s", subMessage, fm.Italic(fm.Link(commit["description"].(string), commit["url"].(string))))

				jiraDomain := os.Getenv("JIRA_DOMAIN")

				if jiraDomain != "" {
					jira := commit["jira"].([]string)

					if len(jira) != 0 {
						var jiraMessage []string
						for _, j := range jira {
							url := fmt.Sprintf("%s/browse/%s", jiraDomain, strings.ToUpper(j))
							jiraMessage = append(jiraMessage, fm.Link(j, url))
						}

						subMessage = fmt.Sprintf("%s (%s)", subMessage, strings.Join(jiraMessage, ", "))
					}
				}
			}
		}
		message = fmt.Sprintf("%s%s", message, subMessage)
	}

	return message
}

func (t *PipelineLogType) Make() string {
	message, err := t.Header()

	if err != nil {
		return ""
	}

	return fmt.Sprintf("%s%s\n%s", message, t.Body(), t.Footer())
}

type MergeDefaultType struct {
	Event     *gitlab.MergeEvent
	Subscribe *models.Subscribe
}

func (t *MergeDefaultType) Make() string {
	var message string
	url := fm.Link(t.Event.Project.Name, t.Event.Project.WebURL)
	if t.Event.ObjectAttributes.MergeStatus == "unchecked" {
		message = fmt.Sprintf("🎭⚠ Необходимо проверить MergeRequest! | %s (%d)", url, t.Event.Project.ID)
	} else if t.Event.ObjectAttributes.MergeStatus == "cannot_be_merged" {
		message = fmt.Sprintf("🎭❌ Обнаружены ошибки в MergeRequest! | %s (%d)", url, t.Event.Project.ID)
	} else if t.Event.ObjectAttributes.MergeStatus == "can_be_merged" {
		message = fmt.Sprintf("🎭✅ Был завершён MergeRequest! | %s (%d)", url, t.Event.Project.ID)
	} else {
		return ""
	}

	message = fmt.Sprintf("%s\n—————\n%s", message, fm.Link(t.Event.ObjectAttributes.Title, t.Event.ObjectAttributes.URL))
	message = fmt.Sprintf("%s\n\n🌳: %s → %s", message, tgbotapi.EscapeText(tgbotapi.ModeHTML, t.Event.ObjectAttributes.SourceBranch), tgbotapi.EscapeText(tgbotapi.ModeHTML, t.Event.ObjectAttributes.TargetBranch))
	url = fmt.Sprintf("%s/%s", os.Getenv("GITLAB_URL"), t.Event.User.Username)
	message = fmt.Sprintf("%s\n🧙: %s", message, fm.Link(t.Event.User.Name, url))

	return message
}
