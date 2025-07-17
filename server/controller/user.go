package controller

import (
	"net/http"
	"strconv"
	"strings"

	"monsoon/api"
	"monsoon/db"
	"monsoon/db/app"
	"monsoon/lib"
	"monsoon/util"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserDB         app.IUserDB
	PasswordHasher lib.IPasswordHasher
	TokenHelper    lib.IJWTTokenHelper
}

// @Summary      Create a new user
// @Description  User creation/registration route
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request  body     	api.UserCreateSchema  true  "User info"
// @Success      201      {object}  api.APIResponse
// @Failure      400      {object}  api.APIResponse
// @Router       /user/create [post]
func (ctrl *UserController) UserCreateRoute(c *gin.Context) {

	// Validating input
	req := &api.UserCreateSchema{}
	err := util.ValidateRequestInput(c, req)

	rs := &api.APIResponse{}
	if err != nil {
		util.WriteAPIError(c, "Invalid input", rs, http.StatusBadRequest)
		return
	}

	// Checks if another user already exists with EITHER the same username OR same email
	fields := map[db.UserColumn]any{
		db.ColUserEmail:    req.Email,
		db.ColUserUsername: req.Username,
	}
	user, err := ctrl.UserDB.GetUserByAnyField(c.Request.Context(), fields)
	if err != nil {
		util.HandleServerError(c, rs, err)
		return
	}

	if user != nil {
		util.WriteAPIError(c, "User already exists", rs, http.StatusConflict)
		return
	}

	userId := lib.GenerateSnowflakeID()
	pwHash, err := ctrl.PasswordHasher.Hash(req.Password)
	if err != nil {
		util.HandleServerError(c, rs, err)
		return
	}

	// Refresh token expires after 7 days, gets stored in DB
	rnd_code, _ := util.RandomBase16String(32)
	refreshToken, err := ctrl.TokenHelper.CreateNewToken(rnd_code, 86400*7)
	if err != nil {
		util.HandleServerError(c, rs, err)
		return
	}

	// Creating the user here
	err = ctrl.UserDB.CreateUser(c.Request.Context(), userId.Int64(), strings.ToLower(req.Username), req.DisplayName, req.Email, pwHash, refreshToken)
	if err != nil {
		util.HandleServerError(c, rs, err)
		return
	}
	c.JSON(http.StatusCreated, rs)
}

// @Summary      Login User
// @Description  User authentication route
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request  body     	api.UserLoginSchema  true  "User credentials"
// @Success      200      {object}  api.APIResponse
// @Failure      401      {object}  api.APIResponse
// @Router       /user/login [post]
func (ctrl *UserController) UserLoginRoute(c *gin.Context) {
	/* User login and authentication route */

	req := &api.UserLoginSchema{}
	err := util.ValidateRequestInput(c, req)
	rs := &api.APIResponse{}
	if err != nil {
		util.WriteAPIError(c, "Invalid input", rs, http.StatusBadRequest)
		return
	}

	// Check for user's existence
	fields := map[db.UserColumn]any{
		db.ColUserEmail: req.Email,
	}
	user, err := ctrl.UserDB.GetUserByAnyField(c.Request.Context(), fields)
	if err != nil {
		util.HandleServerError(c, rs, err)
		return
	}

	// Small helper function to handle invalid credential error
	_handleInvalidCredentials := func() {
		util.WriteAPIError(c, "Invalid credentials", rs, http.StatusUnauthorized)
	}

	if user == nil {
		// Handle for invalid user
		_handleInvalidCredentials()
		return
	}

	// Verify password
	ok, err := ctrl.PasswordHasher.Verify(req.Password, user.Password)
	if err != nil {
		util.HandleServerError(c, rs, err)
		return
	}

	if !ok {
		// If your password is wrong, you can piss off
		_handleInvalidCredentials()
		return
	}

	rs.Data = user

	// Update new refresh token in DB upon login
	code, _ := util.RandomBase16String(32)
	refreshToken, err := ctrl.TokenHelper.CreateNewToken(code, api.EXPIRY_REFRESH_TOKEN)
	if err != nil {
		util.HandleServerError(c, rs, err)
		return
	}

	values := map[db.UserColumn]string{
		db.ColUserRefreshToken: refreshToken,
	}
	id, _ := strconv.ParseInt(user.ID, 10, 64)
	err = ctrl.UserDB.UpdateUserTableById(c, id, db.TableAuth, values)
	if err != nil {
		util.HandleServerError(c, rs, err)
		return
	}

	util.SetRefreshTokenCookie(c, refreshToken)
	c.JSON(http.StatusOK, rs)
}

// @Summary      Get Access Token
// @Description  Retrieve authentication access token for authenticated user
// @Tags         users
// @Produce      json
// @Success      200      {object}  api.APIResponse
// @Failure      401      {object}  api.APIResponse
// @Router       /user/token [post]
func (ctrl *UserController) UserGetAccessToken(c *gin.Context) {
	rs := &api.APIResponse{}

	user, ok := util.GetCurrentUser(c)
	if !ok {
		util.WriteAPIError(c, "Unauthorized", rs, http.StatusUnauthorized)
		return
	}

	var exp int64

	config := util.GetConfig()
	if config.IsDev {
		// 1 hour expiry in dev env, PITA to hit the login route every few minutes
		exp = 60 * 60
	} else {
		exp = api.EXPIRY_ACCESS_TOKEN
	}
	accessToken, err := ctrl.TokenHelper.CreateNewToken(user.ID, exp)
	if err != nil {
		util.HandleServerError(c, rs, err)
		return
	}
	rs.Data = accessToken
	c.JSON(http.StatusOK, rs)
}

// @Summary      Get Current User
// @Description  Currently authenticated user route
// @Tags         users
// @Produce      json
// @Success      200      {object}  api.APIResponse
// @Failure      401      {object}  api.APIResponse
// @Router       /user/@me [post]
func (ctrl *UserController) UserGetCurrent(c *gin.Context) {
	rs := &api.APIResponse{}
	user, ok := util.GetCurrentUser(c)
	if !ok {
		util.WriteAPIError(c, "Unauthorized", rs, http.StatusUnauthorized)
		return
	}
	rs.Data = user
	c.JSON(http.StatusOK, rs)
}
