package middleware

import (
	"net/http"

	"monsoon/lib"

	"github.com/gin-gonic/gin"
)

/* Middleware for routes which require user authentication */

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		rs := &lib.APIResponse{Err: true, Message: "Authorization token required"}
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, rs)
			return
		}

		// Store user ID in context
		// TODO: Use proper user IDs after auth from DB, for now its just using the token
		c.Set("user_id", token)

		// Continue to handler
		c.Next()
	}
}
