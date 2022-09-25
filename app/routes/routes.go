package routes

import (
	"github.com/gorilla/mux"
	"github.com/xanzy/go-gitlab"
	"net/http"
)

func New(Secret string) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", WebIndex).Methods("GET")
	router.HandleFunc("/project/{project_id}/pipeline/{pipeline_id}", WebPipeline).Methods("GET")
	wh := Webhook{
		Secret: Secret,
		EventsToAccept: []gitlab.EventType{
			gitlab.EventTypeMergeRequest,
			gitlab.EventTypePipeline,
			gitlab.EventTypePush,
		},
	}

	router.HandleFunc("/webhook", wh.ServeHTTP).Methods("POST")
	router.Handle("/webhook/toggle/{user_id}", ToggleWebhookSignature(http.HandlerFunc(WebToggle))).Methods("POST")
	router.HandleFunc("/webhook/toggle/{user_id}", GetWebToggle).Methods("GET")

	return router
}
