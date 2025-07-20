package controller

import (
	"monsoon/api"
	"monsoon/db/app"
	"monsoon/util"
	"net/http"

	"github.com/gin-gonic/gin"
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
