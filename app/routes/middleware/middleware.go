package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func PanicRecovery(c *gin.Context, recovered interface{}) {
	if err, ok := recovered.(string); ok {
		c.JSON(
			http.StatusInternalServerError,
			fmt.Sprintf("Error: %s", err),
		)
	}
	c.AbortWithStatus(http.StatusInternalServerError)
}
