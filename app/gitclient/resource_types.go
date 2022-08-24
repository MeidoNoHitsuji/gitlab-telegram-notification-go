package gitclient

import (
	"errors"
	"fmt"
	"github.com/leodido/go-conventionalcommits"
	"github.com/leodido/go-conventionalcommits/parser"
	"github.com/xanzy/go-gitlab"
	"gitlab-telegram-notification-go/helper"
	"gitlab-telegram-notification-go/models"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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
		message = fmt.Sprintf("🧩❌ PipeLine завершился ошибкой. | [%s](%s) (%d)", t.Event.Project.Name, t.Event.Project.WebURL, t.Event.Project.ID)
		message = fmt.Sprintf("%s\n—————", message)
	} else if t.Event.ObjectAttributes.Status == "success" {
		message = fmt.Sprintf("🧩✅ PipeLine завершился успешно. | [%s](%s) (%d)", t.Event.Project.Name, t.Event.Project.WebURL, t.Event.Project.ID)
		message = fmt.Sprintf("%s\n—————", message)
	} else {
		return "", errors.New("Такой статус пайплайна не поддерживается.")
	}

	if t.Event.MergeRequest.ID != 0 {
		message = fmt.Sprintf("%s\n[%s](%s/-/pipelines/%d)\n—————", message, t.Event.MergeRequest.Title, t.Event.Project.WebURL, t.Event.ObjectAttributes.ID)
	} else {
		messages := strings.Split(t.Event.Commit.Message, "\n")

		if len(messages) > 0 {
			message = fmt.Sprintf("%s\n[%s](%s/-/pipelines/%d)\n—————", message, messages[0], t.Event.Project.WebURL, t.Event.ObjectAttributes.ID)
		} else {
			message = fmt.Sprintf("%s\n[%s](%s/-/pipelines/%d)\n—————", message, t.Event.Commit.Message, t.Event.Project.WebURL, t.Event.ObjectAttributes.ID)
		}
	}

	return message, nil
}

func (t PipelineDefaultType) Footer() string {
	var message string
	if t.Event.MergeRequest.ID != 0 {
		message = fmt.Sprintf("\n🌳: %s → %s", t.Event.MergeRequest.SourceBranch, t.Event.MergeRequest.TargetBranch)
	} else {
		message = fmt.Sprintf("\n🌳: %s", t.Event.ObjectAttributes.Ref)
	}

	message = fmt.Sprintf("%s\n🧙: [%s](%s/%s)", message, t.Event.User.Name, os.Getenv("GITLAB_URL"), t.Event.User.Username)

	return message
}

func (t *PipelineDefaultType) Make() string {
	message, err := t.Header()
	if err != nil {
		return ""
	}
	message = fmt.Sprintf("%s\nСборочная линия:", message)

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

				message = fmt.Sprintf("%s\n%s [%s](%s/-/jobs/%d)", message, emoji, build.Name, t.Event.Project.WebURL, build.ID)
			}
		}
	}

	return fmt.Sprintf("%s\n%s", message, t.Footer())
}

type PipelineCommitsType struct {
	PipelineDefaultType
	Commits []*gitlab.Commit
}

func (t *PipelineCommitsType) Make() string {
	message, err := t.Header()
	if err != nil {
		return ""
	}
	message = fmt.Sprintf("%s\nЗалитые изменения:", message)

	for _, commit := range t.Commits {
		if len(commit.ParentIDs) > 1 {
			continue
		}
		message = fmt.Sprintf("%s\n📄 [%s](%s)", message, commit.Title, commit.WebURL)
	}

	return fmt.Sprintf("%s\n%s", message, t.Footer())
}

type PipelineLogType struct {
	PipelineDefaultType
	Commits []*gitlab.Commit
}

func (t *PipelineLogType) Make() string {
	message, err := t.Header()

	if err != nil {
		return ""
	}

	message = fmt.Sprintf("%s\n", message)

	commits := map[string]map[string][]map[string]interface{}{}

	for _, commit := range t.Commits {
		if len(commit.ParentIDs) > 1 {
			continue
		}

		res, _ := parser.NewMachine(parser.WithBestEffort()).Parse([]byte(commit.Message))

		if res.Ok() {
			resCommit := res.(*conventionalcommits.ConventionalCommit)

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

			if scope == nil {
				scope = gitlab.String("Другое")
			}

			_, ok = commits[t][*scope]

			if !ok {
				commits[t][*scope] = []map[string]interface{}{}
			}

			body := resCommit.Body

			if body == nil {
				body = gitlab.String("")
			}

			jira, ok := resCommit.Footers["jira"]

			if !ok {
				jira = []string{}
			}

			commits[t][*scope] = append(commits[t][*scope], map[string]interface{}{
				"description": resCommit.Description,
				"url":         commit.WebURL,
				"body":        *body,
				"jira":        jira,
			})
		}
	}

	for k, v := range types {
		data, ok := commits[k]

		if !ok {
			continue
		}

		subMessage := fmt.Sprintf("*%s*:\n", v)
		for scopeKey, dataCommits := range data {
			subMessage = fmt.Sprintf("%s    __%s__:\n", subMessage, cases.Title(language.Und).String(scopeKey))
			for _, commit := range dataCommits {
				subMessage = fmt.Sprintf("%s        📄_[%s](%s)_", subMessage, commit["description"], commit["url"])

				jiraDomain := os.Getenv("JIRA_DOMAIN")

				if jiraDomain != "" {
					jira := commit["jira"].([]string)

					if len(jira) != 0 {
						var jiraMessage []string
						for _, j := range jira {
							jiraMessage = append(jiraMessage, fmt.Sprintf("[%s](%s/browse/%s)", j, jiraDomain, j))
						}

						//Это какой-то пздц.. Почему в golang нельзя экранировать символы прямо в тексте!??!?!
						subMessage = fmt.Sprintf("%s %s%s%s", subMessage, `\(`, strings.Join(jiraMessage, ", "), `\)`)
					}
				}

				subMessage = fmt.Sprintf("%s\n", subMessage)
			}
		}
		message = fmt.Sprintf("%s%s\n", message, subMessage)
	}

	return fmt.Sprintf("%s\n%s", message, t.Footer())
}

type MergeDefaultType struct {
	Event     *gitlab.MergeEvent
	Subscribe *models.Subscribe
}

func (t *MergeDefaultType) Make() string {
	var message string
	if t.Event.ObjectAttributes.MergeStatus == "unchecked" {
		message = fmt.Sprintf("🎭⚠ Необходимо проверить MergeRequest. | [%s](%s) (%d)", t.Event.Project.Name, t.Event.Project.WebURL, t.Event.Project.ID)
	} else if t.Event.ObjectAttributes.MergeStatus == "cannot_be_merged" {
		message = fmt.Sprintf("🎭❌ Обнаружены ошибки в MergeRequest. | [%s](%s) (%d)", t.Event.Project.Name, t.Event.Project.WebURL, t.Event.Project.ID)
	} else if t.Event.ObjectAttributes.MergeStatus == "can_be_merged" {
		message = fmt.Sprintf("🎭✅ Был завершён MergeRequest. | [%s](%s) (%d)", t.Event.Project.Name, t.Event.Project.WebURL, t.Event.Project.ID)
	} else {
		return ""
	}

	message = fmt.Sprintf("%s\n—————\n[%s](%s)", message, t.Event.ObjectAttributes.Title, t.Event.ObjectAttributes.URL)
	message = fmt.Sprintf("%s\n\n🌳: %s → %s", message, t.Event.ObjectAttributes.SourceBranch, t.Event.ObjectAttributes.TargetBranch)
	message = fmt.Sprintf("%s\n🧙: [%s](%s/%s)", message, t.Event.User.Name, os.Getenv("GITLAB_URL"), t.Event.User.Username)

	return message
}
