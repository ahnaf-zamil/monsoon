package controller

import (
	"errors"
	"log"
	"monsoon/api"
	"monsoon/db/app"
	"monsoon/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type UserController struct {
	UserDB app.IUserDB
}

// @Summary      Get Current User
// @Description  Currently authenticated user route
// @Tags         users
// @Produce      json
// @Success      200      {object}  api.APIResponse
// @Failure      401      {object}  api.APIResponse
// @Router       /user/me [post]
// @Security     BearerAuth
func (ctrl *UserController) UserGetCurrent(c *gin.Context) {
	rs := &api.APIResponse{}
	user, ok := util.GetCurrentUser(c)
	if !ok {
		util.WriteAPIError(c, "Unauthorized", rs, http.StatusUnauthorized)
		return
	}

	util.WriteAPIResponse(c, user, rs, http.StatusOK)
}

// @Summary      Search User
// @Description  Search a user using username
// @Tags         users
// @Produce      json
// @Param        username   path      string  true  "Query username"
// @Success      200      {object}  api.APIResponse
// @Failure      401      {object}  api.APIResponse
// @Router       /user/search/{username} [get]
// @Security     BearerAuth
func (ctrl *UserController) UserSearchUser(c *gin.Context) {
	rs := &api.APIResponse{}
	_, ok := util.GetCurrentUser(c)
	if !ok {
		util.WriteAPIError(c, "Unauthorized", rs, http.StatusUnauthorized)
		return
	}

	queryUsername := c.Param("query")
	users, err := ctrl.UserDB.SearchUsersByUsername(c.Request.Context(), queryUsername)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// In case there are no returned users
			arr := [...]any{}
			util.WriteAPIResponse(c, arr, rs, http.StatusOK)
			return
		}
		log.Println(err)
		util.HandleServerError(c, rs, err)
		return
	}

	util.WriteAPIResponse(c, users, rs, http.StatusOK)
}
