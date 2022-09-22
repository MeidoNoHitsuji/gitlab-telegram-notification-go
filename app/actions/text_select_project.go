package actions

import (
	"encoding/json"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/xanzy/go-gitlab"
	"gitlab-telegram-notification-go/actions/callbacks"
	"gitlab-telegram-notification-go/database"
	"gitlab-telegram-notification-go/gitclient"
	"gitlab-telegram-notification-go/models"
	"gitlab-telegram-notification-go/telegram"
)

const SelectProjectActionType ActionNameType = "sp_act" //select_project

type SelectProjectAction struct {
	BaseAction
	CallbackData *callbacks.ChangeActiveProjectType `json:"callback_data"`
	BackData     SelectProjectBackData              `json:"bd"`
}

type SelectProjectBackData struct {
	ProjectId int `json:"pi"`
}

func (act *SelectProjectAction) SetIsBack(update tgbotapi.Update) error {

	err := act.BaseAction.SetIsBack(update)

	if err != nil {
		return err
	}

	var tmp map[string]interface{}
	err = json.Unmarshal([]byte(update.CallbackQuery.Data), &tmp)

	if err != nil {
		return err
	}

	backData, ok := tmp["bd"]

	if !ok {
		return errors.New("Параметры не найдены.")
	}

	out, _ := json.Marshal(backData)
	err = json.Unmarshal(out, &act.BackData)

	if err != nil {
		return err
	}

	return nil
}

func (act *SelectProjectAction) Validate(update tgbotapi.Update) bool {
	if !act.BaseAction.Validate(update) {
		return false

	}

	if act.InitializedBy == InitByCallback {
		return json.Unmarshal([]byte(update.CallbackQuery.Data), &act.CallbackData) == nil
	} else {
		return true
	}
}

func (act *SelectProjectAction) Active(update tgbotapi.Update) error {
	message, botMessage := telegram.GetMessageFromUpdate(update)

	if message == nil {
		return errors.New("Неизвестно откуда прилетел запрос.")
	}

	var project *gitlab.Project
	var err error
	git := gitclient.Instant()

	if act.IsBack {
		project, _, err = git.Projects.GetProject(act.BackData.ProjectId, &gitlab.GetProjectOptions{})

		if err != nil {
			return err
		}
	} else if act.InitializedBy == InitByCallback {
		project, _, err = git.Projects.GetProject(act.CallbackData.ProjectId, &gitlab.GetProjectOptions{})

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

	db := database.Instant()

	projectObj := models.Project{
		ID: project.ID,
	}

	db.FirstOrCreate(&projectObj, models.Project{
		Name: project.Name,
	})

	subscribeObj := database.FirstOrCreateSubscribe(project.ID, message.Chat.ID, true)

	backData := callbacks.NewBackType(nil)
	backOut, err := json.Marshal(backData)

	if err != nil {
		return err
	}

	settingsData := callbacks.NewSelectProjectSettingsType(project.ID)
	settingsOut, err := json.Marshal(settingsData)

	if err != nil {
		return err
	}

	changeActiveData := callbacks.NewChangeActiveProjectType(project.ID)
	changeActive, err := json.Marshal(changeActiveData)

	if err != nil {
		return err
	}

	var textAct string

	if act.InitializedBy == InitByCallback {
		if !subscribeObj.DeletedAt.Valid {
			db.Delete(&subscribeObj)
			textAct = "Вкл"
		} else {
			db.Exec("UPDATE subscribes SET deleted_at = NULL WHERE id = ?", subscribeObj.ID)
			textAct = "Выкл"
		}
	} else {
		if !subscribeObj.DeletedAt.Valid {
			textAct = "Выкл"
		} else {
			textAct = "Вкл"
		}
	}

	keyboard := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(textAct, string(changeActive)),
		tgbotapi.NewInlineKeyboardButtonData("Настройки", string(settingsOut)),
	)

	keyboardBack := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Отмена", string(backOut)),
	)

	if botMessage {
		telegram.UpdateMessageById(
			message,
			fmt.Sprintf("Был выбран проект: %s", project.Name),
			tgbotapi.NewInlineKeyboardMarkup(keyboard, keyboardBack),
			nil,
		)
	} else {
		telegram.SendMessageById(
			message.Chat.ID,
			fmt.Sprintf("Был выбран проект: %s", project.Name),
			tgbotapi.NewInlineKeyboardMarkup(keyboard, keyboardBack),
			nil,
		)
	}

	return nil
}

func NewSelectProjectAction() *SelectProjectAction {
	return &SelectProjectAction{
		BaseAction: BaseAction{
			ID: SelectProjectActionType,
			InitBy: []ActionInitByType{
				InitByCallback,
				InitByText,
			},
			AfterAction:           SubscribesActionType,
			InitCallbackFuncNames: []callbacks.CallbackFuncName{callbacks.ChangeActiveFuncName},
			BeforeAction:          SubscribesActionType,
		},
	}
}
