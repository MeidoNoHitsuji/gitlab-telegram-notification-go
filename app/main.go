package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/xanzy/go-gitlab"
	"gitlab-telegram-notification-go/client"
	"gitlab-telegram-notification-go/command"
	"gitlab-telegram-notification-go/database"
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
	bot := client.Telegram()

	go runWebServer(os.Getenv("GITLAB_SECRET"), os.Getenv("WEBHOOK_PORT"))

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	chat, err := bot.GetChat(tgbotapi.ChatInfoConfig{
		ChatConfig: tgbotapi.ChatConfig{
			//ChatID: 479413765,
			SuperGroupUsername: "meidonohitsuji",
		},
	})

	if err != nil {
		log.Panic(err)
	}

	msg := tgbotapi.NewMessage(chat.ID, "kekw.")

	if _, err := bot.Send(msg); err != nil {
		log.Panic(err)
	}

	//git.Projects.AddProjectHook()

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message updates
			continue
		}

		if !update.Message.IsCommand() { // ignore any non-command Messages
			continue
		}

		// Create a new MessageConfig. We don't have text yet,
		// so we leave it empty.
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		// Extract the command from the Message.
		switch update.Message.Command() {
		case "subscribe":
			text, _ := command.Subscribe(update.Message.CommandArguments())
			msg.Text = text
		case "start":
			msg.Text = "I'm ok."
		default:
			msg.Text = "I don't know that command"
		}

		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
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
	mux.Handle("/webhook", wh)

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
