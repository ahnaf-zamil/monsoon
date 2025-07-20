package controller

import (
	"bytes"
	"crypto/rand"
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

type AuthController struct {
	UserDB         app.IUserDB
	PasswordHasher lib.IPasswordHasher
	TokenHelper    lib.IJWTTokenHelper
}

// @Summary      Create a new user
// @Description  User creation/registration route
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body     	api.UserCreateSchema  true  "User info"
// @Success      201      {object}  api.APIResponse
// @Failure      400,409      {object}  api.APIResponse
// @Router       /auth/create [post]
func (ctrl *AuthController) AuthRegistrationRoute(c *gin.Context) {
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

	userID := lib.GenerateSnowflakeID()

	// Creating the user here
	err = ctrl.UserDB.CreateUser(c.Request.Context(), userID.Int64(), strings.ToLower(req.Username), req.DisplayName, req.Email, decodedFields["pwHash"], decodedFields["encKey"], decodedFields["sigKey"], decodedFields["authSalt"], decodedFields["encSalt"], decodedFields["encSeed"], decodedFields["nonce"])
	if err != nil {
		util.HandleServerError(c, rs, err)
		return
	}

	// Creating user session (Essentially logging him in right after registration for better UX)
	sessionID := lib.GenerateSnowflakeID()

	// Create refreshToken with session ID
	refreshToken, err := ctrl.TokenHelper.CreateNewToken(sessionID.String(), api.EXPIRY_REFRESH_TOKEN)
	if err != nil {
		util.HandleServerError(c, rs, err)
		return
	}

	// Creating session entry
	err = ctrl.UserDB.CreateUserSession(c.Request.Context(), sessionID.Int64(), userID.Int64(), refreshToken)
	if err != nil {
		util.HandleServerError(c, rs, err)
		return
	}

	util.SetRefreshTokenCookie(c, refreshToken)

	// TODO: Generate public key hashes and store in Merkle tree for later verification

	util.WriteAPIResponse(c, user, rs, http.StatusCreated)
}

// @Summary      User salt
// @Description  Get a user's salt by email
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body     	api.UserLoginSaltSchema  true  "User email"
// @Success      200      {object}  api.APIResponse
// @Failure      400      {object}  api.APIResponse
// @Router       /auth/salt [post]
func (ctrl *AuthController) AuthFetchUserSalt(c *gin.Context) {
	/* Fetch user's password salt */

	req := &api.UserLoginSaltSchema{}
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
	user, _ := ctrl.UserDB.GetUserByAnyField(c.Request.Context(), fields)
	// No need for error checks, auth wont exist if user doesnt exist
	userAuth, _ := ctrl.UserDB.GetUserAuthByID(c.Request.Context(), user.ID)

	var salt []byte
	if userAuth == nil {
		// Return dummy salt to prevent enumeration
		salt = make([]byte, 32) // 32 byte random salt
		_, err := rand.Read(salt)
		if err != nil {
			util.HandleServerError(c, rs, err)
			return
		}
	} else {
		salt = userAuth.PasswordSalt
	}

	util.WriteAPIResponse(c, salt, rs, http.StatusOK)
}

// @Summary      Login user
// @Description  Log into user account
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body     	api.UserLoginSchema  true  "User credentials"
// @Success      200      {object}  api.APIResponse
// @Failure      400      {object}  api.APIResponse
// @Router       /auth/login [post]
func (ctrl *AuthController) AuthLoginUser(c *gin.Context) {
	/* Authenticate user using password hash */

	req := &api.UserLoginSchema{}
	err := util.ValidateRequestInput(c, req)
	rs := &api.APIResponse{}
	if err != nil {
		util.WriteAPIError(c, "Invalid input", rs, http.StatusBadRequest)
		return
	}
	decodedPwHash, err := base64.StdEncoding.DecodeString(req.PasswordHash)
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
		util.WriteAPIError(c, "Unauthorized", rs, http.StatusUnauthorized)
		return
	}

	if !bytes.Equal(decodedPwHash, user.PasswordHash) {
		// If sent password hash and user's password hash in DB are unequal
		util.WriteAPIError(c, "Unauthorized", rs, http.StatusUnauthorized)
		return
	}

	// Creating user session (Essentially logging him in right after registration for better UX)
	sessionID := lib.GenerateSnowflakeID()

	// Create refreshToken with session ID
	refreshToken, err := ctrl.TokenHelper.CreateNewToken(sessionID.String(), api.EXPIRY_REFRESH_TOKEN)
	if err != nil {
		util.HandleServerError(c, rs, err)
		return
	}

	// Creating session entry
	intUserID, err := strconv.ParseInt(user.ID, 10, 64)
	if err != nil {
		util.HandleServerError(c, rs, err)
		return
	}

	err = ctrl.UserDB.CreateUserSession(c.Request.Context(), sessionID.Int64(), intUserID, refreshToken)
	if err != nil {
		util.HandleServerError(c, rs, err)
		return
	}

	util.SetRefreshTokenCookie(c, refreshToken)

	userAuth, _ := ctrl.UserDB.GetUserAuthByID(c.Request.Context(), user.ID)
	util.WriteAPIResponse(c, userAuth, rs, http.StatusOK)
}

// @Summary      Get Access Token
// @Description  Retrieve authentication access token for authenticated user
// @Tags         auth
// @Produce      json
// @Success      200      {object}  api.APIResponse
// @Failure      401      {object}  api.APIResponse
// @Router       /auth/token [post]
func (ctrl *AuthController) UserGetAccessToken(c *gin.Context) {
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
