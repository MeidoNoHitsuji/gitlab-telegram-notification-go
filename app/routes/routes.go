package routes

import (
	"github.com/gin-gonic/gin"
	middleware2 "github.com/s12i/gin-throttle"
	"github.com/xanzy/go-gitlab"
	"gitlab-telegram-notification-go/routes/middleware"
)

func New(Secret string) *gin.Engine {
	router := gin.Default()
	router.Use(gin.CustomRecovery(middleware.PanicRecovery))
	router.LoadHTMLGlob("static/*")

	router.GET("/", WebIndex)
	router.GET("/project/:project_id/pipeline/:pipeline_id", WebPipeline)

	webhookRouter := router.Group("/webhook")
	{
		wh := Webhook{
			Secret: Secret,
			EventsToAccept: []gitlab.EventType{
				gitlab.EventTypeMergeRequest,
				gitlab.EventTypePipeline,
				gitlab.EventTypeNote,
				gitlab.EventTypePush,
			},
		}

		webhookRouter.POST("/", wh.ServeHTTP)

		toggleRouter := webhookRouter.Group("/toggle/:user_telegram_id")
		toggleRouter.Use(middleware2.Throttle(1, 1))
		{
			toggleRouter.POST("/", WebToggle)
		}
	}

	apiRouter := router.Group("/api")
	{
		apiRouter.GET("/test", TestFunc)
	}

	oauthRouter := router.Group("/oauth")
	{
		oauthRouter.GET("/gitlab", GitlabOAuth)
		oauthRouter.GET("/jira", JiraOAuth)
	}

	return router
}
