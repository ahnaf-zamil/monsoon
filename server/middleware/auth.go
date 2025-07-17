package middleware

import (
	"monsoon/api"
	"monsoon/db"
	"monsoon/db/app"
	"monsoon/lib"
	"net/http"

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
		token := c.GetHeader("Authorization")
		rs := &api.APIResponse{Err: true, Message: "Unauthorized"}

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
