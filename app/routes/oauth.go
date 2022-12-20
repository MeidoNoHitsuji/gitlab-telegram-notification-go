package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab-telegram-notification-go/routes/middleware"
	"net/http"
	"os"
)

func GitlabOAuth(c *gin.Context) {
	code := "testcode"
	user := middleware.GetUserByApi(c)

	if user != nil {
		//TODO: Авторизовывать
	} else {
		url := fmt.Sprintf("https://t.me/%s?start=gitlab-%s", os.Getenv("TELEGRAM_USERNAME"), code)
		c.Redirect(http.StatusSeeOther, url)
	}
}

func JiraOAuth(c *gin.Context) {
	code := "testcode"
	user := middleware.GetUserByApi(c)

	if user != nil {
		//TODO: Авторизовывать
	} else {
		url := fmt.Sprintf("https://t.me/%s?start=jira-%s", os.Getenv("TELEGRAM_USERNAME"), code)
		c.Redirect(http.StatusSeeOther, url)
	}
}
