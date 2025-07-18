package controller

import (
	"monsoon/db/app"
	"monsoon/lib"
	"monsoon/middleware"
	"monsoon/ws"

	"github.com/gin-gonic/gin"
)

func InitControllers(r *gin.Engine) {
	/* Registering HTTP route controllers by creating subrouters */

	tokenHelper := lib.GetJWTTokenHelper()
	userDB := app.GetUserDB()
	passHasher := lib.GetPasswordHasher()

	rfTokenMiddleware := middleware.RefreshTokenRequired(userDB, tokenHelper)
	authMiddleware := middleware.AuthRequired(userDB, tokenHelper)

	api := r.Group("/api")
	msg := api.Group("/message")
	user := api.Group("/user")

	msg_ctrl := &MessageController{UserDB: userDB, NATS_PUB: &ws.NATSPublisher{}}
	msg.POST("/user/:recipientID", authMiddleware, msg_ctrl.MessageUserRoute)

	user_ctrl := &UserController{UserDB: userDB, PasswordHasher: passHasher, TokenHelper: tokenHelper}
	user.POST("/create", user_ctrl.UserCreateRoute)
	user.POST("/login", user_ctrl.UserLoginRoute)
	user.POST("/token", rfTokenMiddleware, user_ctrl.UserGetAccessToken)
	user.GET("/me", authMiddleware, user_ctrl.UserGetCurrent)
}
