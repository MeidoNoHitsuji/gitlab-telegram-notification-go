package actions

import (
	"encoding/json"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/xanzy/go-gitlab"
	"gitlab-telegram-notification-go/actions/callbacks"
	"gitlab-telegram-notification-go/gitclient"
	"gitlab-telegram-notification-go/telegram"
)

const SelectProjectActionType ActionNameType = "select_project"

type SelectProjectAction struct {
	BaseAction
}

type SelectProjectBackData struct {
	projectId int
}

func (act *SelectProjectAction) Active(update tgbotapi.Update) error {
	message, _ := telegram.GetMessageFromUpdate(update)

	if message == nil {
		return errors.New("Неизвестно откуда прилетел запрос.")
	}

	var project *gitlab.Project
	var err error
	git := gitclient.Instant()

	if act.IsBack {
		backData := act.BackData.(SelectProjectBackData)
		project, _, err = git.Projects.GetProject(backData.projectId, &gitlab.GetProjectOptions{})

		if err != nil {
			return err
		}
	} else {
		projects, _, err := git.Projects.ListProjects(&gitlab.ListProjectsOptions{
			Search:           gitlab.String(message.Text),
			SearchNamespaces: gitlab.Bool(true),
		})

		if err != nil {
			return err
		}

		if len(projects) == 0 {
			return NewErrorForUser("Не было найдено ни единого проекта по вашему запросу.")
		}

		project = projects[0]
	}

	backData := callbacks.NewBackType()
	backOut, err := json.Marshal(backData)

	if err != nil {
		return err
	}

	settingsData := callbacks.NewSelectProjectSettingsType()

	settingsOut, err := json.Marshal(settingsData)

	if err != nil {
		return err
	}

	fmt.Println(string(settingsOut))

	keyboard := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Вкл/Выкл", "{}"),
		tgbotapi.NewInlineKeyboardButtonData("Настройки", string(settingsOut)),
	)

	keyboardBack := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Отмена", string(backOut)),
	)

	telegram.SendMessageById(
		message.Chat.ID,
		fmt.Sprintf("Был выбран проект: %s", project.Name),
		tgbotapi.NewInlineKeyboardMarkup(keyboard, keyboardBack),
		nil,
	)

	return nil
}

func NewSelectProjectAction() *SelectProjectAction {
	return &SelectProjectAction{
		BaseAction: BaseAction{
			ID:          SelectProjectActionType,
			InitBy:      InitByText,
			AfterAction: SubscribesActionType,
		},
	}
}
