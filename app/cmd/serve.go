package cmd

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/cobra"
	"gitlab-telegram-notification-go/actions"
	"gitlab-telegram-notification-go/database"
	"gitlab-telegram-notification-go/routes"
	"gitlab-telegram-notification-go/telegram"
	"log"
	"net/http"
	"os"
	"time"
)

var serveCmd = &cobra.Command{
	Use: "serve",
	Run: serve,
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func serve(cmd *cobra.Command, args []string) {

	database.Instant()
	bot := telegram.Instant()

	go runWebServer(os.Getenv("GITLAB_SECRET"), os.Getenv("RUN_PORT"))

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {

		if update.MyChatMember != nil {
			if update.MyChatMember.Chat.Type == "private" {
				if update.MyChatMember.OldChatMember.Status == "kicked" && update.MyChatMember.NewChatMember.Status == "member" {
					database.UpdateMemberStatus(update.MyChatMember.Chat.ID, update.MyChatMember.From.UserName, false)
				} else if update.MyChatMember.OldChatMember.Status == "member" && update.MyChatMember.NewChatMember.Status == "kicked" {
					database.UpdateMemberStatus(update.MyChatMember.Chat.ID, update.MyChatMember.From.UserName, true)
				}
			} else if update.MyChatMember.Chat.Type == "group" {
				if update.MyChatMember.OldChatMember.Status == "left" && update.MyChatMember.NewChatMember.Status == "member" {
					database.UpdateChatStatus(update.MyChatMember.Chat.ID, false)
				} else if update.MyChatMember.OldChatMember.Status == "member" && update.MyChatMember.NewChatMember.Status == "left" {
					database.UpdateChatStatus(update.MyChatMember.Chat.ID, true)
				}
			}
		}

		if update.Message != nil {
			database.UpdateMemberStatus(update.Message.Chat.ID, update.Message.From.UserName, false)
		}

		err := actions.Active(update)

		if err != nil {
			switch err := err.(type) {
			case actions.ErrorForUser:
				if update.CallbackQuery != nil {
					callback := tgbotapi.NewCallbackWithAlert(
						update.CallbackQuery.ID,
						err.Error(),
					)
					bot.Request(callback)
				} else {
					telegram.SendMessageById(update.Message.Chat.ID, err.Error(), nil, nil)
				}
			default:
				fmt.Println(err.Error())
			}
		} else {
			if update.CallbackQuery != nil {
				callback := tgbotapi.NewCallback(
					update.CallbackQuery.ID,
					"",
				)
				bot.Request(callback)
			}
		}
	}
}

func runWebServer(Secret string, Port string) {

	srv := &http.Server{
		Handler: routes.New(Secret),
		Addr:    fmt.Sprintf("0.0.0.0:%s", Port),

		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
}
