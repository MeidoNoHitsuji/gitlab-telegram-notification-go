package routes

import (
	"github.com/gin-gonic/gin"
	middleware2 "github.com/s12i/gin-throttle"
	"github.com/xanzy/go-gitlab"
	"gitlab-telegram-notification-go/routes/middleware"
	"time"
)

func New(Secret string) *gin.Engine {
	router := gin.Default()

	router.Use(gin.CustomRecovery(middleware.PanicRecovery))

	router.LoadHTMLGlob("static/*")
	router.GET("/", WebIndex)
	router.GET("/project/:project_id/pipeline/:pipeline_id", WebPipeline)

	wh := Webhook{
		Secret: Secret,
		EventsToAccept: []gitlab.EventType{
			gitlab.EventTypeMergeRequest,
			gitlab.EventTypePipeline,
			gitlab.EventTypeNote,
			gitlab.EventTypePush,
		},
	}

	router.POST("/webhook", wh.ServeHTTP)

	webhookToggle := router.Group("/webhook/toggle/:user_telegram_id")
	webhookToggle.Use(middleware2.Throttle(1000, 1))
	{
		webhookToggle.POST("/", WebToggle)
	}

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
