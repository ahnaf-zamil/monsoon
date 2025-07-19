package controller

import (
	"encoding/base64"
	"monsoon/api"
	"monsoon/db"
	"monsoon/db/app"
	"monsoon/lib"
	"monsoon/util"
	"net/http"
	"strconv"
	"strings"

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
// @Failure      400,409      {object}  api.APIResponse
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

	b64Fields := map[string]string{
		"encKey":   req.Keys.Enc,
		"sigKey":   req.Keys.Sig,
		"pwHash":   req.PasswordHash,
		"encSeed":  req.EncryptedSeed,
		"nonce":    req.Nonce,
		"encSalt":  req.Salts.Enc,
		"authSalt": req.Salts.Auth,
	}

	// Validate encryption and signature key base64
	decodedFields := make(map[string][]byte)
	for fieldName, b64 := range b64Fields {
		decoded, err := base64.StdEncoding.DecodeString(b64)
		if err != nil {
			util.WriteAPIError(c, "Invalid input", rs, http.StatusBadRequest)
			return
		}
		decodedFields[fieldName] = decoded
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

	// Creating the user here
	err = ctrl.UserDB.CreateUser(c.Request.Context(), userId.Int64(), strings.ToLower(req.Username), req.DisplayName, req.Email, decodedFields["pwHash"], decodedFields["encKey"], decodedFields["sigKey"], decodedFields["authSalt"], decodedFields["encSalt"], decodedFields["encSeed"], decodedFields["nonce"])
	if err != nil {
		util.HandleServerError(c, rs, err)
		return
	}

	// TODO: Generate public key hashes and store in Merkle tree for later verification

	util.WriteAPIResponse(c, user, rs, http.StatusCreated)
}

// @Summary      Login User
// @Description  User authentication route
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request  body     	api.UserLoginSchema  true  "User credentials"
// @Success      200      {object}  api.APIResponse
// @Failure      400,401      {object}  api.APIResponse
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
	// ok, err := ctrl.PasswordHasher.Verify(req.Password, user.Password)
	// if err != nil {
	// 	util.HandleServerError(c, rs, err)
	// 	return
	// }

	// if !ok {
	// 	// If your password is wrong, you can piss off
	// 	_handleInvalidCredentials()
	// 	return
	// }

	// Update new refresh token in DB upon login
	code, _ := util.RandomBase16String(32)
	refreshToken, err := ctrl.TokenHelper.CreateNewToken(code, api.EXPIRY_REFRESH_TOKEN)
	if err != nil {
		util.HandleServerError(c, rs, err)
		return
	}

	userID, _ := strconv.ParseInt(user.ID, 10, 64)
	sessionID := lib.GenerateSnowflakeID()

	// Creating session entry
	err = ctrl.UserDB.CreateUserSession(c.Request.Context(), sessionID.Int64(), userID, refreshToken)
	if err != nil {
		util.HandleServerError(c, rs, err)
		return
	}

	util.SetRefreshTokenCookie(c, refreshToken)
	util.WriteAPIResponse(c, user, rs, http.StatusOK)
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

		// 2 hour expiry in dev env, PITA to hit the token route every few minutes
		exp = 2 * 60 * 60
	} else {
		exp = api.EXPIRY_ACCESS_TOKEN
	}
	accessToken, err := ctrl.TokenHelper.CreateNewToken(user.ID, exp)
	if err != nil {
		util.HandleServerError(c, rs, err)
		return
	}
	util.WriteAPIResponse(c, accessToken, rs, http.StatusOK)
}

// @Summary      Get Current User
// @Description  Currently authenticated user route
// @Tags         users
// @Produce      json
// @Success      200      {object}  api.APIResponse
// @Failure      401      {object}  api.APIResponse
// @Router       /user/me [post]
// @Security    BearerAuth
func (ctrl *UserController) UserGetCurrent(c *gin.Context) {
	rs := &api.APIResponse{}
	user, ok := util.GetCurrentUser(c)
	if !ok {
		util.WriteAPIError(c, "Unauthorized", rs, http.StatusUnauthorized)
		return
	}

	util.WriteAPIResponse(c, user, rs, http.StatusOK)
}
