package jiraclient

import (
	"fmt"
	"github.com/andygrunwald/go-jira"
	"github.com/trivago/tgo/tmath"
	"gitlab-telegram-notification-go/database"
	"gitlab-telegram-notification-go/models"
	"gitlab-telegram-notification-go/routes/request"
	"strconv"
	"strings"
)

func UpdateJiraWorklog(telegramChannelId int64, data request.ToggleData) {

	if data.Payload.Duration < 0 {
		return
	}

	issueKey := GetIssueKeyFromText(data.Payload.Description)

	if issueKey == "" {
		fmt.Println("Ключ не найден!!")
		return
	}

	client, err := New(telegramChannelId)

	if err != nil {
		fmt.Println("Ошибка создания слиента")
		fmt.Println(err)
		return
	}

	issue, _, err := client.Issue.Get(issueKey, &jira.GetQueryOptions{})

	if err != nil {
		fmt.Println("Ошибка запроса задачи " + issueKey)
		fmt.Println(err)
		return
	}

	if strings.ToLower(issue.Fields.Status.Name) == "закрыто" {
		fmt.Println("Нельзя редактировать закрытые задачи")
		return
	}

	me, _, err := client.User.GetSelf()

	if err != nil {
		fmt.Println("Ошибка получения самого себя")
		fmt.Println(err)
		return
	}

	if issue.Fields.Assignee.AccountID != me.AccountID {
		fmt.Println("Попытка отредактировать задачу, которая не принадлежит тебе")
		return
	}

	IssueId, err := strconv.Atoi(issue.ID)

	if err != nil {
		fmt.Println("Ошибка парсинга issue.ID - " + issue.ID)
		fmt.Println(err)
		return
	}

	db := database.Instant()
	var eventIntegration models.ToggleJiraIntegration

	res := db.Where(models.ToggleJiraIntegration{
		TimeEntityId: data.Payload.Id,
	}).Find(&eventIntegration)

	if res.RowsAffected == 0 {
		eventIntegration.TimeEntityId = data.Payload.Id
		eventIntegration.IssueId = IssueId

		worklogRecord, _, err := client.Issue.AddWorklogRecord(issue.ID, &jira.WorklogRecord{
			Comment:          "Created by TelegramBot",
			Started:          getTime(data.Payload.Start),
			TimeSpentSeconds: tmath.MaxI(data.Payload.Duration, 60),
		})

		if err != nil {
			fmt.Println("Ошибка добавления времени")
			fmt.Println(err)
			return
		}

		eventIntegration.WorklogRecordId, err = strconv.Atoi(worklogRecord.ID)

		if err != nil {
			fmt.Println("Ошибка парсинга worklogRecord.ID - " + worklogRecord.ID)
			fmt.Println(err)
			return
		}

		db.Save(&eventIntegration)
	} else {
		if eventIntegration.IssueId != IssueId {
			fmt.Println("Пересоздаём карточку") //Мне лень расписывать отдельно по запросам
			DeleteJiraWorklog(telegramChannelId, data)
			UpdateJiraWorklog(telegramChannelId, data)
		} else {
			IssueId := strconv.Itoa(eventIntegration.IssueId)
			WorklogId := strconv.Itoa(eventIntegration.WorklogRecordId)

			_, _, err := client.Issue.UpdateWorklogRecord(IssueId, WorklogId, &jira.WorklogRecord{
				Started:          getTime(data.Payload.Start),
				TimeSpentSeconds: tmath.MaxI(data.Payload.Duration, 60),
			})

			if err != nil {
				fmt.Println("Ошибка обновлнения времени")
				fmt.Println(err)
				return
			} else {
				fmt.Println("Время обновлено!")
				return
			}
		}
	}
}

func DeleteJiraWorklog(telegramChannelId int64, data request.ToggleData) {
	issueKey := GetIssueKeyFromText(data.Payload.Description)

	if issueKey == "" {
		fmt.Println("Ключ не найден!!")
		return
	}

	client, err := New(telegramChannelId)

	if err != nil {
		fmt.Println("Ошибка создания клиента")
		fmt.Println(err)
		return
	}

	db := database.Instant()
	var eventIntegration models.ToggleJiraIntegration

	res := db.Where(models.ToggleJiraIntegration{
		TimeEntityId: data.Payload.Id,
	}).Find(&eventIntegration)

	if res.RowsAffected == 0 {
		return
	}

	IssueId := strconv.Itoa(eventIntegration.IssueId)
	WorklogId := strconv.Itoa(eventIntegration.WorklogRecordId)

	worklogs, _, err := client.Issue.GetWorklogs(IssueId)

	if err != nil {
		fmt.Println("Ошибка получения ворклогов " + issueKey)
		return
	}

	found := false

	for _, worklog := range worklogs.Worklogs {
		if worklog.ID == WorklogId {
			found = true
		}
	}

	if !found {
		db.Delete(&eventIntegration)
		fmt.Println("Ворклог не был найден в jira и был удалён!")
		return
	}

	issue, _, err := client.Issue.Get(IssueId, &jira.GetQueryOptions{})

	if err != nil {
		fmt.Println("Ошибка запроса задачи " + issueKey)
		fmt.Println(err)
		return
	}

	if strings.ToLower(issue.Fields.Status.Name) == "закрыто" {
		fmt.Println("Нельзя редактировать закрытые задачи")
		return
	}

	me, _, err := client.User.GetSelf()

	if err != nil {
		fmt.Println("Ошибка получения самого себя")
		fmt.Println(err)
		return
	}

	if issue.Fields.Assignee.AccountID != me.AccountID {
		fmt.Println("Попытка отредактировать задачу, которая не принадлежит тебе")
		return
	}

	_, err = DeleteWorklogRecord(client, IssueId, WorklogId)
	if err != nil {
		fmt.Println("Ворклог не удалось удалить!")
		fmt.Println(err)
		return
	}

	db.Delete(&eventIntegration)
	fmt.Println("Ворклог удалён!")
}
