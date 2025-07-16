package controller

import (
	"monsoon/db/app"
	"monsoon/lib"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserDB         app.IUserDB
	PasswordHasher lib.IPasswordHasher
}

// @Summary      Create a new user
// @Description  Register a new user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request  body     	lib.UserCreateSchema  true  "User info"
// @Success      201      {object}  lib.APIResponse
// @Failure      400      {object}  lib.APIResponse
// @Router       /user/create [post]
func (ctrl *UserController) UserCreateRoute(c *gin.Context) {
	/* User creation/registration route */

	// Validating input
	req := &lib.UserCreateSchema{}
	err := lib.ValidateRequestInput(c, req)

	rs := &lib.APIResponse{}
	if err != nil {
		lib.WriteAPIError(c, "Invalid input", rs, http.StatusBadRequest)
		return
	}

	// Checks if another user already exists with EITHER the same username OR same email
	fields := map[lib.UserColumn]any{
		lib.ColUserEmail:    req.Email,
		lib.ColUserUsername: req.Username,
	}
	user, err := ctrl.UserDB.GetUserByAnyField(c.Request.Context(), fields)
	if err != nil {
		lib.HandleServerError(c, rs, err)
		return
	}

	if user != nil {
		lib.WriteAPIError(c, "User already exists", rs, http.StatusConflict)
		return
	}

	user_id := lib.GenerateSnowflakeID()
	pw_hash, err := ctrl.PasswordHasher.Hash(req.Password)
	if err != nil {
		lib.HandleServerError(c, rs, err)
		return
	}

	// Creating the user here
	err = ctrl.UserDB.CreateUser(c.Request.Context(), user_id.Int64(), strings.ToLower(req.Username), req.DisplayName, req.Email, pw_hash)
	if err != nil {
		lib.HandleServerError(c, rs, err)
		return
	}
	c.JSON(http.StatusCreated, rs)
}

// @Summary      Login User
// @Description  User authentication route
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request  body     	lib.UserLoginSchema  true  "User credentials"
// @Success      200      {object}  lib.APIResponse
// @Failure      401      {object}  lib.APIResponse
// @Router       /user/login [post]
func (ctrl *UserController) UserLoginRoute(c *gin.Context) {
	/* User login and authentication route */

	req := &lib.UserLoginSchema{}
	err := lib.ValidateRequestInput(c, req)
	rs := &lib.APIResponse{}
	if err != nil {
		lib.WriteAPIError(c, "Invalid input", rs, http.StatusBadRequest)
		return
	}

	// Check for user's existence
	fields := map[lib.UserColumn]any{
		lib.ColUserEmail: req.Email,
	}
	user, err := ctrl.UserDB.GetUserByAnyField(c.Request.Context(), fields)
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
	ok, err := ctrl.PasswordHasher.Verify(req.Password, user.Password)
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
