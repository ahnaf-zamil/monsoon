package controller

import (
	"log"
	"net/http"
	"ws_realtime_app/db/app"
	"ws_realtime_app/lib"

	"github.com/gin-gonic/gin"
	"github.com/matthewhartstonge/argon2"
)

func UserCreateRoute(c *gin.Context) {
	// Validating input
	req := &lib.UserCreateSchema{}
	rs, err := lib.ValidateRequestInput(c, req)
	if err != nil {
		return
	}

	argon := argon2.DefaultConfig()

	user_id := lib.GenerateSnowflakeID()
	pw_hash, err := argon.HashEncoded([]byte(req.Password))
	if err != nil {
		lib.WriteAPIError(c, "Internal server error", rs, http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// TODO: Check if user with email/username already exists and handle gracefully

	err = app.Users.CreateUser(c.Request.Context(), user_id.Int64(), req.Username, req.DisplayName, req.Email, pw_hash)
	if err != nil {
		lib.WriteAPIError(c, "Internal server error", rs, http.StatusInternalServerError)
		log.Println(err)
		return
	}
	c.JSON(http.StatusCreated, rs)
}
