package controller

import "github.com/gin-gonic/gin"

// import "github.com/gorilla/mux"
func InitControllers(r *gin.Engine) {
	/* Registering HTTP route controllers by creating subrouters */

	api := r.Group("/api")
	msg := api.Group("/message")

	msg.POST("/create/:room_id", MessageCreateRoute)
}
