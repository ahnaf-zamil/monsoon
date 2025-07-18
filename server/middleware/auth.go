package middleware

import (
	"monsoon/api"
	"monsoon/db"
	"monsoon/db/app"
	"monsoon/lib"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

/* Middleware for routes which require user authentication */

func RefreshTokenRequired(userDB app.IUserDB, jwt lib.IJWTTokenHelper) gin.HandlerFunc {
	/* Middleware to check whether user has refresh token cookie
	to get access token */
	return func(c *gin.Context) {
		token, err := c.Cookie(api.COOKIE_REFRESH_TOKEN)

		rs := &api.APIResponse{Err: true, Message: "Unauthorized"}
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, rs)
			return
		}

		_, err = jwt.VerifyToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, rs)
			return
		}

		// Fetch user from DB
		fields := map[db.UserColumn]any{
			db.ColUserRefreshToken: token,
		}
		user, err := userDB.GetUserByAnyField(c.Request.Context(), fields)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, rs)
			return
		}

		// Store user in context
		c.Set("current_user", user)

		// Continue to handler
		c.Next()
	}
}

func AuthRequired(userDB app.IUserDB, jwt lib.IJWTTokenHelper) gin.HandlerFunc {
	/* Middleware to authenticate user using short-lived access token
	sent in Authorization header */

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		rs := &api.APIResponse{Err: true, Message: "Unauthorized"}

		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, rs)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, rs)
			return
		}

		token := strings.TrimSpace(parts[1])
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, rs)
			return
		}

		userID, err := jwt.VerifyToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, rs)
			return
		}

		fields := map[db.UserColumn]any{
			db.ColUserID: userID,
		}
		user, err := userDB.GetUserByAnyField(c.Request.Context(), fields)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, rs)
			return
		}

		// Store user in context
		c.Set("current_user", user)

		// Continue to handler
		c.Next()
	}
}
