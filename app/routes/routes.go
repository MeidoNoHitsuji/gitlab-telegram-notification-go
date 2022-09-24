package routes

import (
	"github.com/gorilla/mux"
	"github.com/xanzy/go-gitlab"
	"gitlab-telegram-notification-go/webhook"
)

func New(Secret string) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", WebIndex).Methods("GET")
	router.HandleFunc("/project/{project_id}/pipeline/{pipeline_id}", WebPipeline).Methods("GET")
	wh := webhook.Webhook{
		Secret: Secret,
		EventsToAccept: []gitlab.EventType{
			gitlab.EventTypeMergeRequest,
			gitlab.EventTypePipeline,
			gitlab.EventTypePush,
		},
	}

	router.HandleFunc("/webhook", wh.ServeHTTP).Methods("POST")
	router.HandleFunc("/webhook/toggle/{user_id}", WebToggle).Methods("POST")

	return router
}
