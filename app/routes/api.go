package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func TestFunc(c *gin.Context) {
	c.JSON(
		http.StatusOK,
		gin.H{
			"message": "Hello Api!",
		},
	)
}
