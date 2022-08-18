package database

import (
	"errors"
	"fmt"
	"github.com/xanzy/go-gitlab"
	"gitlab-telegram-notification-go/helper"
	"gitlab-telegram-notification-go/models"
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

	builder := db.Model(&models.Subscribe{}).Preload("TelegramChannel").Joins("inner join subscribe_events as event on event.subscribe_id = subscribes.id")

	if filter.ProjectId != 0 {
		builder = builder.Where("subscribes.project_id = ?", filter.ProjectId)
	}

	if filter.Event != "" {
		builder = builder.Where("event.event = ?", filter.Event)
	}

	if filter.Status != "" {
		p := "JSON_EXTRACT(event.parameters, '$[*].status')"
		builder = builder.Where(fmt.Sprintf("(%s is null or %s = '' or %s = ?)", p, p, p), filter.Status)
	}

	if filter.AuthorUsername != "" {
		p := "JSON_EXTRACT(event.parameters, '$[*].author_username')"
		builder = builder.Where(fmt.Sprintf("(%s is null or %s = '' or %s = ?)", p, p, p), filter.AuthorUsername)
	}

	if filter.BranchName != "" {
		p := "JSON_EXTRACT(event.parameters, '$[*].branch_name')"
		builder = builder.Where(fmt.Sprintf("(%s is null or %s = '' or %s = ?)", p, p, p), filter.BranchName)
	}

	builder = builder.Find(&subscribes)

	return subscribes
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
