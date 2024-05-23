package ginmw

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/ahaostudy/calendar_reminder/utils/jwt"
)

func GlobalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		token = strings.TrimPrefix(token, "Bearer ")

		if id, ok := jwt.ParseToken(token); ok {
			c.Set("user_id", id)
		}
		c.Next()
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, exists := c.Get("user_id"); !exists {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Next()
	}
}
