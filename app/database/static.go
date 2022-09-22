package database

import (
	"errors"
	"fmt"
	"github.com/xanzy/go-gitlab"
	"gitlab-telegram-notification-go/helper"
	"gitlab-telegram-notification-go/models"
	"gorm.io/gorm"
	"strings"
)

func UpdateMemberStatus(telegramId int64, username string, isDeleted bool) *models.User {
	channel := UpdateChatStatus(telegramId, isDeleted)

	db := Instant()
	user := models.User{
		TelegramChannelId: channel.ID,
	}

	db.Model(&models.User{}).FirstOrCreate(&user)
	username = strings.ToLower(username)
	if user.Username != username {
		user.Username = username
		db.Save(&user)
	}

	return &user
}

func UpdateChatStatus(telegramId int64, isDeleted bool) *models.TelegramChannel {
	db := Instant()
	channel := models.TelegramChannel{
		ID: telegramId,
	}
	db.Model(&models.TelegramChannel{}).FirstOrCreate(&channel)
	if channel.Active != !isDeleted {
		channel.Active = !isDeleted
		db.Save(&channel)
	}
	return &channel
}

func UpdateSubscribes(project gitlab.Project, telegramId int64, events ...string) error {
	db := Instant()

	events = helper.Unique(events)

	projectObj := models.Project{
		ID: project.ID,
	}

	db.Model(models.Project{}).FirstOrCreate(&projectObj)

	if projectObj.Name != project.Name {
		projectObj.Name = project.Name
		db.Save(&projectObj)
	}

	telegram := models.TelegramChannel{}
	result := db.Where(models.TelegramChannel{ID: telegramId}).Find(&telegram)

	if result.RowsAffected == 0 {
		return errors.New(fmt.Sprintf("Канал с Id %d не был найден.", telegramId))
	}

	subscribe := models.Subscribe{
		ProjectId:         project.ID,
		TelegramChannelId: telegramId,
	}

	db.Where(&subscribe).Preload("Events").FirstOrCreate(&subscribe)

	for _, event := range subscribe.Events {
		if !helper.Contains(events, event.Event) {
			db.Where(&models.SubscribeEvent{
				SubscribeId: subscribe.ID,
				Event:       event.Event,
			}).Delete(&models.SubscribeEvent{})
		}
	}

	for _, event := range events {
		db.Where(&models.SubscribeEvent{
			SubscribeId: subscribe.ID,
			Event:       event,
		}).FirstOrCreate(&models.SubscribeEvent{})
	}

	return nil
}

func GetSubscribesByProjectIdAndKind(filter GetSubscribesFilter) []models.Subscribe {
	var subscribes []models.Subscribe
	db := Instant()

	builder := db.Model(&models.Subscribe{}).Preload("TelegramChannel").Preload("Events").Joins("inner join subscribe_events as event on event.subscribe_id = subscribes.id")

	if filter.ProjectId != 0 {
		builder = builder.Where("subscribes.project_id = ?", filter.ProjectId)
	}

	if filter.Event != "" {
		builder = builder.Where("event.event = ?", filter.Event)
	}

	if filter.Status != "" {
		p1 := "JSON_EXTRACT(event.parameters, '$.status') is null"
		p2 := "JSON_LENGTH(JSON_EXTRACT(event.parameters, '$.status')) = 0"
		p3 := "JSON_CONTAINS(event.parameters, JSON_ARRAY(?), '$.status')"
		builder = builder.Where(fmt.Sprintf("(%s or %s or %s)", p1, p2, p3), filter.Status)
	}

	if filter.AuthorUsername != "" {
		p1 := "JSON_EXTRACT(event.parameters, '$.author_username') is null"
		p2 := "JSON_LENGTH(JSON_EXTRACT(event.parameters, '$.author_username')) = 0"
		p3 := "JSON_CONTAINS(event.parameters, JSON_ARRAY(?), '$.author_username')"
		builder = builder.Where(fmt.Sprintf("(%s or %s or %s)", p1, p2, p3), filter.AuthorUsername)
	}

	if filter.ToBranchName != "" {
		p1 := "JSON_EXTRACT(event.parameters, '$.to_branch_name') is null"
		p2 := "JSON_LENGTH(JSON_EXTRACT(event.parameters, '$.to_branch_name')) = 0"
		p3 := "JSON_CONTAINS(event.parameters, JSON_ARRAY(?), '$.to_branch_name')"
		builder = builder.Where(fmt.Sprintf("(%s or %s or %s)", p1, p2, p3), filter.ToBranchName)
	}

	if filter.FromBranchName != "" {
		p1 := "JSON_EXTRACT(event.parameters, '$.to_branch_name') is null"
		p2 := "JSON_LENGTH(JSON_EXTRACT(event.parameters, '$.to_branch_name')) = 0"
		p3 := "JSON_CONTAINS(event.parameters, JSON_ARRAY(?), '$.to_branch_name')"
		builder = builder.Where(fmt.Sprintf("(%s or %s or %s)", p1, p2, p3), filter.FromBranchName)
	}

	if filter.IsMerge != "" {
		p1 := "JSON_EXTRACT(event.parameters, '$.to_branch_name') is null"
		p2 := "JSON_LENGTH(JSON_EXTRACT(event.parameters, '$.to_branch_name')) = 0"
		p3 := "JSON_CONTAINS(event.parameters, JSON_ARRAY(?), '$.to_branch_name')"
		builder = builder.Where(fmt.Sprintf("(%s or %s or %s)", p1, p2, p3), filter.IsMerge)
	}

	builder = builder.Find(&subscribes)

	return subscribes
}

func GetProjectsByTelegramIdsWithDeleted(telegramIds ...int64) []models.Project {
	var projects []models.Project
	db := Instant()

	builder := db.Unscoped().Model(&models.Project{}).Joins("inner join subscribes as subscribe on subscribe.project_id = projects.id")
	builder = builder.Where("subscribe.telegram_channel_id in ?", telegramIds)
	builder = builder.Group("projects.id").Limit(9)
	builder = builder.Find(&projects)

	return projects
}

func GetProjectsByTelegramIds(telegramIds ...int64) []models.Project {
	var projects []models.Project
	db := Instant()

	builder := db.Model(&models.Project{}).Joins("inner join subscribes as subscribe on subscribe.project_id = projects.id")
	builder = builder.Where("subscribe.telegram_channel_id in ?", telegramIds)
	builder = builder.Group("projects.id").Limit(9)
	builder = builder.Find(&projects)

	return projects
}

func GetEventsByProjectId(projectId int) []string {
	var subscribes []models.Subscribe
	var events []string
	db := Instant()

	builder := db.Model(&models.Subscribe{}).Preload("Events").Where("subscribes.project_id = ?", projectId)
	builder = builder.Find(&subscribes)

	for _, subscribe := range subscribes {
		for _, event := range subscribe.Events {
			events = append(events, event.Event)
		}
	}

	return helper.Unique(events)
}

func GetUserActionInChannel(telegramId int64, username string) *models.UserTelegramChannelAction {
	db := Instant()

	obj := models.UserTelegramChannelAction{}

	builder := db.Model(&models.UserTelegramChannelAction{}).Joins("inner join users as user on user_telegram_channel_actions.user_id = user.id")
	builder = builder.Where("user_telegram_channel_actions.telegram_channel_id = ?", telegramId)
	builder = builder.Where("user.username = ?", username)
	builder = builder.Find(&obj)

	if builder.RowsAffected == 0 {
		return nil
	} else {
		return &obj
	}
}

func UpdateUserActionInChannel(telegramId int64, username string, action string) error {
	obj := GetUserActionInChannel(telegramId, username)

	if obj == nil {
		if err := CreateUserActionInChannel(telegramId, username, action); err != nil {
			return err
		}
	} else {
		db := Instant()
		db.Model(&models.UserTelegramChannelAction{}).Where(obj).Update("action", action)
	}

	return nil
}

func CreateUserActionInChannel(telegramId int64, username string, action string) error {
	db := Instant()

	obj := models.UserTelegramChannelAction{}

	user := &models.User{
		Username: username,
	}

	r := db.Model(&models.User{}).Find(&user)

	if r.RowsAffected == 0 {
		return errors.New("Такой пользователь не найден!")
	}

	obj.Action = action
	obj.TelegramChannelId = telegramId
	obj.UserId = user.ID
	obj.Parameters = ""
	db.Create(obj)

	return nil
}

func UpdateUserActionParameterInChannel(telegramId int64, username string, parameters string) error {
	obj := GetUserActionInChannel(telegramId, username)

	if obj == nil {
		return errors.New("Action не был найден!")
	} else {
		db := Instant()
		db.Model(&models.UserTelegramChannelAction{}).Where(obj).Update("parameters", parameters)
	}

	return nil
}

func UpdateUserActionFormatterInChannel(telegramId int64, username string, formatter string) error {
	obj := GetUserActionInChannel(telegramId, username)

	if obj == nil {
		return errors.New("Action не был найден!")
	} else {
		db := Instant()
		db.Model(&models.UserTelegramChannelAction{}).Where(obj).Update("formatter", formatter)
	}

	return nil
}

func FirstOrCreateSubscribe(ProjectId int, TelegramChannelId int64, WithDelete bool) *models.Subscribe {
	db := Instant()

	subscribeObj := models.Subscribe{
		ProjectId:         ProjectId,
		TelegramChannelId: TelegramChannelId,
	}

	var builder *gorm.DB

	if WithDelete {
		builder = db.Unscoped().Model(&models.Subscribe{})
	} else {
		builder = db.Model(&models.Subscribe{})
	}

	res := builder.Where(&subscribeObj).Find(&subscribeObj)

	if WithDelete {
		builder = db.Unscoped().Model(&models.Subscribe{})
	} else {
		builder = db.Model(&models.Subscribe{})
	}

	if res.RowsAffected == 0 {
		builder.Save(&subscribeObj)
	}

	return &subscribeObj
}
