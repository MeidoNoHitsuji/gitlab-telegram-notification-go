package routes

import (
	"github.com/gorilla/mux"
	"github.com/xanzy/go-gitlab"
	"gitlab-telegram-notification-go/webhook"
)

func New(Secret string) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", WebIndex).Methods("GET")

	wh := webhook.Webhook{
		Secret: Secret,
		EventsToAccept: []gitlab.EventType{
			gitlab.EventTypeMergeRequest,
			gitlab.EventTypePipeline,
			gitlab.EventTypePush,
		},
	}

	router.HandleFunc("/webhook", wh.ServeHTTP).Methods("POST")

	return router
}
