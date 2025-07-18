package controller

import (
	"monsoon/api"
	"monsoon/db/app"
	"monsoon/lib"
	"monsoon/util"
	"monsoon/ws"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type MessageController struct {
	NATS_PUB ws.INATSPublisher
	UserDB   app.IUserDB
}

// @Summary      Directly Message a User
// @Description  Send a direct message to a user
// @Tags         messages
// @Accept       json
// @Produce      json
// @Param        userId   path      int  true  "User ID"
// @Param        request  body     	api.MessageCreateSchema  true  "Message data"
// @Success      201      {object}  api.APIResponse
// @Failure      401      {object}  api.APIResponse
// @Router       /message/user/{recipientId} [post]
// @Security    BearerAuth
func (ctrl *MessageController) MessageUserRoute(c *gin.Context) {
	// Validating input
	req := &api.MessageCreateSchema{}
	err := util.ValidateRequestInput(c, req)

	rs := &api.APIResponse{}
	if err != nil {
		util.WriteAPIError(c, "Invalid input", rs, http.StatusBadRequest)
		return
	}
	author, ok := util.GetCurrentUser(c)
	if !ok {
		util.WriteAPIError(c, "Unauthorized", rs, http.StatusUnauthorized)
		return
	}

	recipientID := c.Param("recipientID")

	recipient, err := ctrl.UserDB.GetUserByID(c.Request.Context(), recipientID)
	if err != nil {
		util.WriteAPIError(c, "User not found", rs, http.StatusNotFound)
		return
	}

	content := req.Content
	payload := api.MessageModel{
		ID:          lib.GenerateSnowflakeID().String(),
		Content:     content,
		CreatedAt:   time.Now().Unix(),
		RoomID:      "", // Empty for DM messages
		AuthorID:    author.ID,
		RecipientID: recipient.ID,
		IsDM:        true,
	}

	rs.Data = &payload
	// Dispatch new message to NATS
	ctrl.NATS_PUB.SendMsgNATS(payload)

	// TODO: Persist message in DB

	c.JSON(http.StatusCreated, rs)
}
