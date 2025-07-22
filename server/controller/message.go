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
	NATS_PUB       ws.INATSPublisher
	UserDB         app.IUserDB
	ConversationDB app.IConversationDB
}

// @Summary      Directly Message a User
// @Description  Send a direct message to a user
// @Tags         messages
// @Accept       json
// @Produce      json
// @Param        userId   path      int  true  "Recipient ID"
// @Param        request  body     	api.MessageCreateSchema  true  "Message data"
// @Success      201      {object}  api.APIResponse
// @Failure      401,403,400      {object}  api.APIResponse
// @Router       /message/user/{userId} [post]
// @Security    BearerAuth
func (ctrl *MessageController) MessageUserRoute(c *gin.Context) {
	// Validating input
	req := &api.MessageCreateSchema{}
	err := util.ValidateRequestInput(c, req)
	rs := &api.APIResponse{}

	author, ok := util.GetCurrentUser(c)
	if !ok {
		util.WriteAPIError(c, "Unauthorized", rs, http.StatusUnauthorized)
		return
	}

	if err != nil {
		util.WriteAPIError(c, "Invalid input", rs, http.StatusBadRequest)
		return
	}

	recipientID := c.Param("recipientID")
	recipient, err := ctrl.UserDB.GetUserByID(c.Request.Context(), recipientID)
	if recipient == nil {
		util.WriteAPIError(c, "User not found", rs, http.StatusNotFound)
		return
	}
	if err != nil {
		util.HandleServerError(c, rs, err)
		return
	}

	dm, err := ctrl.ConversationDB.GetExistingDM(c.Request.Context(), author.ID, recipientID)
	if err != nil {
		util.HandleServerError(c, rs, err)
		return
	}

	if author.ID == recipient.ID {
		// Why is bro DMing himself? Does he not have friends?
		util.WriteAPIError(c, "Bro why are you DMing yourself?", rs, http.StatusForbidden)
		return
	}

	var conversationIDStr string
	if dm == nil {
		conversationID := lib.GenerateSnowflakeID()
		err = ctrl.ConversationDB.CreateUserDM(c.Request.Context(), conversationID.Int64(), author.ID, recipient.ID)
		if err != nil {
			util.HandleServerError(c, rs, err)
			return
		}
		conversationIDStr = conversationID.String()
	} else {
		conversationIDStr = dm.ConversationID
	}

	// TODO: Revamp msg model when we set up msg DB
	payload := api.MessageModel{
		ID:             lib.GenerateSnowflakeID().String(),
		Content:        req.Content,
		CreatedAt:      time.Now().Unix(),
		AuthorID:       author.ID,
		ConversationID: conversationIDStr,
	}

	// TODO: Persist msg in DB

	// Dispatch new message to NATS
	ctrl.NATS_PUB.SendMsgNATS(payload)

	util.WriteAPIResponse(c, payload, rs, http.StatusCreated)
}
