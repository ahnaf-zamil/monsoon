package controller

import (
	"errors"
	"log"
	"monsoon/api"
	"monsoon/db/app"
	msg "monsoon/db/message"
	"monsoon/lib"
	"monsoon/util"
	"monsoon/ws"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type MessageController struct {
	NATS_PUB       ws.INATSPublisher
	UserDB         app.IUserDB
	MsgDB          msg.IMessageDB
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
func (ctrl *MessageController) CreateMessageUserRoute(c *gin.Context) {
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

	msgID := lib.GenerateSnowflakeID()
	// TODO: Revamp msg model when we set up msg DB
	payload := api.MessageModel{
		ID:             msgID.String(),
		Content:        req.Content,
		CreatedAt:      time.Now().Unix(),
		AuthorID:       author.ID,
		ConversationID: conversationIDStr,
	}

	err = ctrl.MsgDB.CreateMessage(c.Request.Context(), msgID.Int64(), conversationIDStr, author.ID, req.Content)
	if err != nil {
		log.Println(err)
		util.HandleServerError(c, rs, err)
		return
	}

	// Dispatch new message to NATS
	ctrl.NATS_PUB.SendMsgNATS(payload)

	util.WriteAPIResponse(c, payload, rs, http.StatusCreated)
}

// @Summary      Send Message to a Conversation
// @Description  What the title says
// @Tags         messages
// @Accept       json
// @Produce      json
// @Param        conversationId   path      int  true  "Conversation ID"
// @Param        request  body     	api.MessageCreateSchema  true  "Message data"
// @Success      201      {object}  api.APIResponse
// @Failure      401,403,400,404      {object}  api.APIResponse
// @Router       /message/conversation/{conversationId} [post]
// @Security    BearerAuth
func (ctrl *MessageController) CreateMessageConversationRoute(c *gin.Context) {
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

	conversationID := c.Param("conversationID")

	convo, err := ctrl.ConversationDB.GetUserConversationByID(c.Request.Context(), conversationID, author.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			util.WriteAPIError(c, "Conversation not found", rs, http.StatusNotFound)
			return
		}
		log.Println(err)
		util.HandleServerError(c, rs, err)
		return
	}

	msgID := lib.GenerateSnowflakeID()
	// TODO: Revamp msg model when we set up msg DB
	payload := api.MessageModel{
		ID:             msgID.String(),
		Content:        req.Content,
		CreatedAt:      time.Now().Unix(),
		AuthorID:       author.ID,
		ConversationID: convo.ConversationID,
	}

	// Persist msg in DB
	err = ctrl.MsgDB.CreateMessage(c.Request.Context(), msgID.Int64(), convo.ConversationID, author.ID, req.Content)
	if err != nil {
		log.Println(err)
		util.HandleServerError(c, rs, err)
		return
	}

	// Dispatch new message to NATS
	ctrl.NATS_PUB.SendMsgNATS(payload)

	util.WriteAPIResponse(c, payload, rs, http.StatusCreated)
}

// @Summary      Get Messages
// @Description  Get messages in a conversation
// @Tags         messages
// @Accept       json
// @Produce      json
// @Param        conversationId   path      int  true  "Conversation ID"
// @Success      201      {object}  api.APIResponse
// @Failure      401,404      {object}  api.APIResponse
// @Router       /message/conversation/{conversationId} [get]
// @Security    BearerAuth
func (ctrl *MessageController) GetMessageConversationRoute(c *gin.Context) {
	rs := &api.APIResponse{}

	author, ok := util.GetCurrentUser(c)
	if !ok {
		util.WriteAPIError(c, "Unauthorized", rs, http.StatusUnauthorized)
		return
	}

	conversationID := c.Param("conversationID")

	convo, err := ctrl.ConversationDB.GetUserConversationByID(c.Request.Context(), conversationID, author.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			util.WriteAPIError(c, "Conversation not found", rs, http.StatusNotFound)
			return
		}
		log.Println(err)
		util.HandleServerError(c, rs, err)
		return
	}

	countStr := c.DefaultQuery("count", "10")
	count, err := strconv.Atoi(countStr)
	if err != nil {
		util.WriteAPIError(c, "Invalid 'count' query", rs, http.StatusBadRequest)
	}

	messages, err := ctrl.MsgDB.GetConversationMessages(c.Request.Context(), convo.ConversationID, count)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// In case there are no messages in this convo yet
			arr := [...]any{}
			util.WriteAPIResponse(c, arr, rs, http.StatusOK)
			return
		}
		log.Println(err)
		util.HandleServerError(c, rs, err)
		return
	}

	util.WriteAPIResponse(c, messages, rs, http.StatusCreated)

}
