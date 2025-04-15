package controller

import (
	"ws_realtime_app/middleware"

	"github.com/gin-gonic/gin"
)

// import "github.com/gorilla/mux"
func InitControllers(r *gin.Engine) {
	/* Registering HTTP route controllers by creating subrouters */

	api := r.Group("/api")
	msg := api.Group("/message", middleware.RequireAuth())

	msg.POST("/create/:room_id", MessageCreateRoute)
}
