package cmd

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/cobra"
	"gitlab-telegram-notification-go/command"
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

	_ = database.Instant()
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

			if update.Message.IsCommand() {
				switch update.Message.Command() {
				case "subscribe":
					text, _, err := command.Subscribe(update.Message.Chat.ID, update.Message.CommandArguments())
					if err == nil {
						telegram.SendMessageById(update.Message.Chat.ID, text, nil, nil)
					} else {
						telegram.SendMessageById(update.Message.Chat.ID, fmt.Sprintf("Ошибка! Не удалось подписаться по причине: %s", err), nil, nil)
					}
					break
				case "start":
					telegram.SendMessageById(update.Message.Chat.ID, "Привет! Мой список команд доступен тебе через <code>/</code>.", nil, nil)
					break
				case "test":
					//tgbotapi.
					//ids := []int64{update.Message.Chat.ID}
					//if update.Message.From.ID != update.Message.Chat.ID {
					//	ids = append(ids, update.Message.From.ID)
					//}
					//command.Test(ids...)

					//keyboard := tgbotapi.NewInlineKeyboardRow(
					//	tgbotapi.InlineKeyboardButton{
					//		Text: "WebApp?",
					//		WebApp: &tgbotapi.WebAppInfo{
					//			URL: "https://gitlab-cicd-tgbot.atwinta.online/",
					//		},
					//	},
					//	tgbotapi.InlineKeyboardButton{
					//		Text: "Pipeline?",
					//		WebApp: &tgbotapi.WebAppInfo{
					//			URL: "https://gitlab-cicd-tgbot.atwinta.online/project/338/pipeline/23885",
					//		},
					//	},
					//)
					//
					//telegram.SendMessageById(update.Message.Chat.ID, "qweqwe", tgbotapi.NewInlineKeyboardMarkup(keyboard), nil)
					break
				case "say":
					deleteMessageConfig := tgbotapi.NewDeleteMessage(update.Message.Chat.ID, update.Message.MessageID)

					_, err := bot.Request(deleteMessageConfig)

					if err != nil {
						fmt.Println(err)
					}

					telegram.SendMessageById(update.Message.Chat.ID, update.Message.CommandArguments(), nil, update.Message.Entities)
					break
				default:
					telegram.SendMessageById(update.Message.Chat.ID, "Я не понимаю, что ты от меня хочешь.", nil, nil)
					break
				}
			} else {
				switch update.Message.Text {
				default:
					//if update.Message.ReplyToMessage != nil {
					//	fmt.Println(1, update.Message.ReplyToMessage, update.Message.ReplyToMessage.ReplyMarkup, update.Message.ReplyMarkup, 2)
					//} else {
					//	fmt.Println(3, update.Message.ReplyMarkup, 4)
					//}
					break
				}
			}

		} else if update.CallbackQuery != nil {

			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			if _, err := bot.Request(callback); err != nil {
				panic(err)
			}

			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
			if _, err := bot.Send(msg); err != nil {
				panic(err)
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
