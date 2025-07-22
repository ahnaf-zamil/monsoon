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
	passHasher := lib.GetPasswordHasher()

	userDB := app.GetUserDB()
	convoDB := app.GetConversationDB()

	rfTokenMiddleware := middleware.RefreshTokenRequired(userDB, tokenHelper)
	authMiddleware := middleware.AuthRequired(userDB, tokenHelper)

	api := r.Group("/api")

	msg := api.Group("/message")
	user := api.Group("/user")
	auth := api.Group("/auth")

	msg_ctrl := &MessageController{UserDB: userDB, NATS_PUB: &ws.NATSPublisher{}, ConversationDB: convoDB}
	msg.POST("/user/:recipientID", authMiddleware, msg_ctrl.MessageUserRoute)

	auth_ctrl := &AuthController{UserDB: userDB, PasswordHasher: passHasher, TokenHelper: tokenHelper}
	auth.POST("/create", auth_ctrl.AuthRegistrationRoute)
	auth.POST("/login", auth_ctrl.AuthLoginUser)
	auth.POST("/salt", auth_ctrl.AuthFetchUserSalt)
	auth.POST("/token", rfTokenMiddleware, auth_ctrl.UserGetAccessToken)

	user_ctrl := &UserController{UserDB: userDB}
	user.GET("/me", authMiddleware, user_ctrl.UserGetCurrent)
}
