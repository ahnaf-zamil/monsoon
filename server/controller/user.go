package controller

import (
	"net/http"
	"strconv"
	"strings"

	"monsoon/db"
	"monsoon/db/app"
	"monsoon/lib"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserDB         app.IUserDB
	PasswordHasher lib.IPasswordHasher
	TokenHelper    lib.IJWTTokenHelper
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
	fields := map[db.UserColumn]any{
		db.ColUserEmail:    req.Email,
		db.ColUserUsername: req.Username,
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

	userId := lib.GenerateSnowflakeID()
	pwHash, err := ctrl.PasswordHasher.Hash(req.Password)
	if err != nil {
		lib.HandleServerError(c, rs, err)
		return
	}

	// Refresh token expires after 7 days, gets stored in DB
	rnd_code, _ := lib.RandomBase16String(32)
	refreshToken, err := ctrl.TokenHelper.CreateNewToken(rnd_code, 86400*7)
	if err != nil {
		lib.HandleServerError(c, rs, err)
		return
	}

	// Creating the user here
	err = ctrl.UserDB.CreateUser(c.Request.Context(), userId.Int64(), strings.ToLower(req.Username), req.DisplayName, req.Email, pwHash, refreshToken)
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
	fields := map[db.UserColumn]any{
		db.ColUserEmail: req.Email,
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

	// Update new refresh token in DB upon login
	rnd_code, _ := lib.RandomBase16String(32)
	refreshToken, err := ctrl.TokenHelper.CreateNewToken(rnd_code, 86400*7)
	if err != nil {
		lib.HandleServerError(c, rs, err)
		return
	}

	values := map[db.UserColumn]string{
		db.ColUserRefreshToken: refreshToken,
	}
	id, _ := strconv.ParseInt(user.ID, 10, 64)
	err = ctrl.UserDB.UpdateUserTableById(c, id, db.TableAuth, values)
	if err != nil {
		lib.HandleServerError(c, rs, err)
		return
	}

	// TODO: Set refresh token in cookie
	lib.SetRefreshTokenCookie(c, refreshToken)
	c.JSON(http.StatusOK, rs)
}
