package actions

import (
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/xanzy/go-gitlab"
	"gitlab-telegram-notification-go/database"
	"gitlab-telegram-notification-go/gitclient"
	"gitlab-telegram-notification-go/helper"
	"gitlab-telegram-notification-go/telegram"
	"log"
	"strings"
)

const SubscribeActionType ActionNameType = "subscribe"

type SubscribeAction struct {
	BaseAction
}

func (act *SubscribeAction) Active(update tgbotapi.Update) error {
	text, _, err := Subscribe(update.Message.Chat.ID, update.Message.CommandArguments())
	if err != nil {
		return NewErrorForUser(fmt.Sprintf("Ошибка! Не удалось подписаться по причине: %s", err))
	} else {
		telegram.SendMessageById(update.Message.Chat.ID, text, nil, nil)
		return nil
	}
}

func NewSubscribeAction() *SelectProjectAction {
	return &SelectProjectAction{
		BaseAction: BaseAction{
			ID:       SubscribeActionType,
			InitBy:   InitByCommand,
			InitText: "subscribe",
		},
	}
}

func getProjectFromArguments(arguments string) (*gitlab.Project, []string, error) {
	git := gitclient.Instant()

	args := strings.Split(
		strings.TrimSpace(arguments), " ")

	if len(args) == 0 || len(args[0]) == 0 {
		return nil, nil, errors.New("Вы должны передать параметром наименование проекта.")
	}

	projects, _, err := git.Projects.ListProjects(&gitlab.ListProjectsOptions{
		Search:           gitlab.String(args[0]),
		SearchNamespaces: gitlab.Bool(true),
	})

	if err != nil {
		log.Print(err)
		return nil, nil, err
	}

	if len(projects) == 0 {
		return nil, nil, errors.New("Не было найдено ни единого проекта по вашему запросу.")
	}

	return projects[0], args[1:], nil
}

func Subscribe(telegramId int64, arguments string) (string, *gitlab.Project, error) {

	project, args, err := getProjectFromArguments(arguments)

	if err != nil {
		return "", nil, err
	}

	var allowArgs []string
	allowEvents := helper.AllowEvents()

	for _, a := range args {
		if len(a) != 0 && helper.Contains(allowEvents, a) {
			allowArgs = append(allowArgs, a)
		}
	}

	if err := database.UpdateSubscribes(*project, telegramId, allowArgs...); err != nil {
		return "", nil, err
	}

	allEvents := database.GetEventsByProjectId(project.ID)

	// TODO: Добавить новые хуки и вынести их в клавиатуру
	text, err := gitclient.Subscribe(project, gitlab.AddProjectHookOptions{
		PushEvents:          gitlab.Bool(helper.Contains(allEvents, helper.Slugify(string(gitlab.EventTypePush)))),
		PipelineEvents:      gitlab.Bool(helper.Contains(allEvents, helper.Slugify(string(gitlab.EventTypePipeline)))),
		MergeRequestsEvents: gitlab.Bool(helper.Contains(allEvents, helper.Slugify(string(gitlab.EventTypeMergeRequest)))),
	})
	if err != nil {
		return "", nil, err
	}

	return text, project, nil
}
