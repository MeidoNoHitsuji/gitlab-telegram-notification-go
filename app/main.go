package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/xanzy/go-gitlab"
	"gitlab-telegram-notification-go/command"
	"gitlab-telegram-notification-go/database"
	"gitlab-telegram-notification-go/telegram"
	"gitlab-telegram-notification-go/webhook"
	"log"
	"net/http"
	"os"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Print(".env not load!!")
	}
}

func main() {
	_ = database.Instant()
	bot := telegram.Instant()

	go runWebServer(os.Getenv("GITLAB_SECRET"), os.Getenv("WEBHOOK_PORT"))

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
						telegram.SendMessageById(update.Message.Chat.ID, text, nil)
					} else {
						telegram.SendMessageById(update.Message.Chat.ID, fmt.Sprintf("Ошибка! Не удалось подписаться по причине: %s", err), nil)
					}
					break
				case "start":
					telegram.SendMessageById(update.Message.Chat.ID, "Привет! Мой список команд доступен тебе через `/`.", nil)
					break
				case "test":
					ids := []int64{update.Message.Chat.ID}
					if update.Message.From.ID != update.Message.Chat.ID {
						ids = append(ids, update.Message.From.ID)
					}
					command.Test(ids...)
					break
				default:
					telegram.SendMessageById(update.Message.Chat.ID, "Я не понимаю, что ты от меня хочешь.", nil)
					break
				}
			} else {
				switch update.Message.Text {
				default:
					if update.Message.ReplyToMessage != nil {
						fmt.Println(1, update.Message.ReplyToMessage, update.Message.ReplyToMessage.ReplyMarkup, update.Message.ReplyMarkup, 2)
					} else {
						fmt.Println(3, update.Message.ReplyMarkup, 4)
					}
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
	wh := webhook.Webhook{
		Secret: Secret,
		EventsToAccept: []gitlab.EventType{
			gitlab.EventTypeMergeRequest,
			gitlab.EventTypePipeline,
			gitlab.EventTypePush,
		},
	}

	mux := http.NewServeMux()
	mux.Handle(fmt.Sprintf("/%s", os.Getenv("WEBHOOK_URL")), wh)

	if err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", Port), mux); err != nil {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
}

//func loop() {
//	msg := tgbotapi.NewMessage(479413765, "Привет")
//	for {
//		if _, err := bot.Send(msg); err != nil {
//			log.Panic(err)
//		}
//		time.Sleep(60 * time.Second)
//	}
//}
