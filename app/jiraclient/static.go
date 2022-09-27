package jiraclient

import (
	"fmt"
	"gitlab-telegram-notification-go/routes/request"
)

func UpdateJiraWorklog(telegramChannelId int64, data request.ToggleData) {
	issueKey := GetIssueKeyFromText(data.Payload.Description)

	if issueKey == "" {
		fmt.Println("Ключ не найден!!")
		return
	}

	_, err := New(telegramChannelId)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("User founded!!")
}

func DeleteJiraWorklog(telegramChannelId int64, data request.ToggleData) {

}
