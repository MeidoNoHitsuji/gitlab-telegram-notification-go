package routes

import (
	"github.com/gorilla/mux"
	"github.com/xanzy/go-gitlab"
	"gitlab-telegram-notification-go/routes/middleware"
	"time"
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
			gitlab.EventTypeNote,
			gitlab.EventTypePush,
		},
	}

	router.HandleFunc("/webhook", wh.ServeHTTP).Methods("POST")
	router.HandleFunc("/webhook/toggle/{user_telegram_id}", WebToggle).Methods("POST")
	router.HandleFunc("/webhook/toggle/{user_id}", GetWebToggle).Methods("GET")
	router.HandleFunc("/panic_test", GetPanic).Methods("GET")

	router.Use(middleware.PanicRecovery)

	go throttleToggleEvents()

	return router
}

func throttleToggleEvents() {
	for {
		newLimited := make(map[int64]time.Time)
		for e, t := range limited {
			if time.Now().Before(t) {
				newLimited[e] = t
			}
		}
		limited = newLimited
		time.Sleep(1 * time.Second)
	}
}
