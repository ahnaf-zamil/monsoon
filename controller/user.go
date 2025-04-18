package controller

import (
	"net/http"
	"strings"
	"ws_realtime_app/db/app"
	"ws_realtime_app/lib"

	"github.com/gin-gonic/gin"
	"github.com/matthewhartstonge/argon2"
)

func UserCreateRoute(c *gin.Context) {
	/* User creation/registration route */

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
		lib.HandleServerError(c, rs, err)
		return
	}

	// TODO: Check if user with email/username already exists and handle gracefully
	// The GetUserByField function will be replaced by GetUserByAnyField() to check for matching email OR username
	user, err := app.GetUserByField(c.Request.Context(), lib.ColUserEmail, req.Email)
	if err != nil {
		lib.HandleServerError(c, rs, err)
		return
	}

	if user != nil {
		lib.WriteAPIError(c, "User already exists", rs, http.StatusConflict)
		return
	}

	err = app.Users.CreateUser(c.Request.Context(), user_id.Int64(), req.Username, strings.ToLower(req.DisplayName), req.Email, pw_hash)
	if err != nil {
		lib.HandleServerError(c, rs, err)
		return
	}
	c.JSON(http.StatusCreated, rs)
}

func UserLoginRoute(c *gin.Context) {
	/* User login and authentication route */

	req := &lib.UserLoginSchema{}
	rs, err := lib.ValidateRequestInput(c, req)
	if err != nil {
		lib.WriteAPIError(c, "Invalid input", rs, http.StatusBadRequest)
		return
	}

	// Check for user's existence
	user, err := app.GetUserByField(c.Request.Context(), lib.ColUserEmail, req.Email)
	if err != nil {
		lib.HandleServerError(c, rs, err)
		return
	}

	// Small helper function to handle invalid credential error
	_handleInvalidCredentials := func() {
		lib.WriteAPIError(c, "Invalid credentials", rs, http.StatusUnauthorized)
	}

	if user == nil {
		// Handle for invalid user
		_handleInvalidCredentials()
		return
	}

	// Verify password
	ok, err := argon2.VerifyEncoded([]byte(req.Password), []byte(user.Password))
	if err != nil {
		lib.HandleServerError(c, rs, err)
		return
	}

	if !ok {
		// If your password is wrong, you can piss off
		_handleInvalidCredentials()
		return
	}

	rs.Data = user

	// TODO: Implement JWT token for auth persistence

	c.JSON(http.StatusOK, rs)
}
