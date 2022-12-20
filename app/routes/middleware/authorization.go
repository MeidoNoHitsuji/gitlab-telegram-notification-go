package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab-telegram-notification-go/configs"
	"gitlab-telegram-notification-go/database"
	"gitlab-telegram-notification-go/models"
)

const UserKey = "user-key"

func GetUserByApi(c *gin.Context) *models.User {

	if user, exists := c.Get(UserKey); exists {
		switch u := user.(type) {
		case *models.User:
			return u
		default:
			delete(c.Keys, UserKey)
		}
	} else {
		if token, err := c.Cookie(configs.CookieAuthToken); err == nil {
			db := database.Instant()
			personal := models.PersonalAccessToken{
				Token: token,
			}

			result := db.Where(&personal).Preload("User")
			if result.RowsAffected != 0 {
				c.Set(UserKey, &personal.User)
				return &personal.User
			}
		} else {
			fmt.Printf("Ошибка авторизации: %s", err.Error())
		}
	}

	return nil
}

func Guest() gin.HandlerFunc {
	return func(c *gin.Context) {
		if GetUserByApi(c) == nil {
			c.Next()
		} else {
			return
			//TODO: Вернуть куда-нибудь
		}
	}
}

func Authorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		if GetUserByApi(c) != nil {
			c.Next()
		} else {
			return
			//TODO: Вернуть на index
		}
	}
}
