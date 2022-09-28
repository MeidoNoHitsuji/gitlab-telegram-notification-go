package jiraclient

import (
	"fmt"
	"github.com/andygrunwald/go-jira"
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

	db := database.Instant()
	var eventIntegration models.ToggleJiraIntegration

	res := db.Where(models.ToggleJiraIntegration{
		TimeEntityId: data.Payload.Id,
	}).Find(&eventIntegration)

	if res.RowsAffected == 0 {
		eventIntegration.TimeEntityId = data.Payload.Id
		eventIntegration.IssueId, err = strconv.Atoi(issue.ID)

		if err != nil {
			fmt.Println("Ошибка парсинга issue.ID - " + issue.ID)
			fmt.Println(err)
			return
		}

		//min := data.Payload.Duration / 60
		//if min == 0 {
		//	min++
		//}

		worklogRecord, _, err := client.Issue.AddWorklogRecord(issue.ID, &jira.WorklogRecord{
			Comment:          "Created by TelegramBot",
			Started:          getTime(data.Payload.Start),
			TimeSpentSeconds: data.Payload.Duration,
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

		fmt.Println("IT'S WORK!!")
	} else {

	}

}

func DeleteJiraWorklog(telegramChannelId int64, data request.ToggleData) {
	//issueKey := GetIssueKeyFromText(data.Payload.Description)
	//
	//if issueKey == "" {
	//	fmt.Println("Ключ не найден!!")
	//	return
	//}
	//
	//client, err := New(telegramChannelId)
	//
	//if err != nil {
	//	fmt.Println("Ошибка создания слиента")
	//	fmt.Println(err)
	//	return
	//}
	//
	//db := database.Instant()
	//var eventIntegration models.ToggleJiraIntegration
	//
	//res := db.Where(models.ToggleJiraIntegration{
	//	TimeEntityId: data.Payload.Id,
	//}).Find(&eventIntegration)
	//
	//if res.RowsAffected == 0 {
	//	return
	//}
	//
	//issue, _, err := client.Issue.Get(fmt.Sprintf("%d", eventIntegration.IssueId), &jira.GetQueryOptions{})
	//
	//if err != nil {
	//	fmt.Println("Ошибка запроса задачи " + issueKey)
	//	fmt.Println(err)
	//	return
	//}
	//
	//if strings.ToLower(issue.Fields.Status.Name) == "закрыто" {
	//	fmt.Println("Нельзя редактировать закрытые задачи")
	//	return
	//}
	//
	//me, _, err := client.User.GetSelf()
	//
	//if err != nil {
	//	fmt.Println("Ошибка получения самого себя")
	//	fmt.Println(err)
	//	return
	//}
	//
	//if issue.Fields.Assignee.AccountID != me.AccountID {
	//	fmt.Println("Попытка отредактировать задачу, которая не принадлежит тебе")
	//	return
	//}
	//
	//TODO: Придётся удалять через вызов request (https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-issue-worklogs/#api-rest-api-3-issue-issueidorkey-worklog-id-delete)
}
