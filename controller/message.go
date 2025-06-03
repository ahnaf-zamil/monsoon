package controller

import (
	"net/http"
	"time"
	"ws_realtime_app/lib"
	"ws_realtime_app/ws"

	"github.com/gin-gonic/gin"
)

type MessageController struct {
	NATS_PUB ws.INATSPublisher
}

// @Summary      Create Message
// @Description  Post a message
// @Tags         messages
// @Accept       json
// @Produce      json
// @Param        roomId   path      int  true  "Room ID"
// @Param        request  body     	lib.MessageCreateSchema  true  "Message data"
// @Success      201      {object}  lib.APIResponse
// @Failure      401      {object}  lib.APIResponse
// @Router       /message/create/{roomId} [post]
// @Security    BearerAuth
func (ctrl *MessageController) MessageCreateRoute(c *gin.Context) {
	// Validating input
	req := &lib.MessageCreateSchema{}
	err := lib.ValidateRequestInput(c, req)

	rs := &lib.APIResponse{}
	if err != nil {
		lib.WriteAPIError(c, "Invalid input", rs, http.StatusBadRequest)
		return
	}

	user_id, exists := c.Get("user_id")
	if !exists {
		lib.WriteAPIError(c, "Auth context error", rs, http.StatusUnauthorized)
		return
	}

	room_id := c.Param("room_id")
	content := req.Content

	// TODO: Use snowflake ID and implement proper payload structuring
	payload := lib.MessageModel{
		ID:        lib.GenerateSnowflakeID().String(),
		Content:   content,
		CreatedAt: time.Now().Unix(),
		RoomID:    room_id,
		UserID:    user_id.(string),
	}

	rs.Data = &payload

	// Dispatch new message to NATS
	// TODO: Dispatch message to Kafka logs for batch processing
	ctrl.NATS_PUB.SendMsgNATS(payload)

	c.JSON(http.StatusCreated, rs)
}
