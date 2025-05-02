package controller

import (
	"ws_realtime_app/db/app"
	"ws_realtime_app/lib"
	"ws_realtime_app/middleware"
	"ws_realtime_app/ws"

	"github.com/gin-gonic/gin"
)

func InitControllers(r *gin.Engine) {
	/* Registering HTTP route controllers by creating subrouters */

	api := r.Group("/api")
	msg := api.Group("/message", middleware.RequireAuth())
	user := api.Group("/user")

	msg_ctrl := &MessageController{NATS_PUB: &ws.NATSPublisher{}}
	msg.POST("/create/:room_id", msg_ctrl.MessageCreateRoute)

	user_ctrl := &UserController{UserDB: app.GetUserDB(), PasswordHasher: lib.GetPasswordHasher()}
	user.POST("/create", user_ctrl.UserCreateRoute)
	user.POST("/login", user_ctrl.UserLoginRoute)
}
