package controller

import (
	"net/http"
	"ws_realtime_app/lib"

	"github.com/gin-gonic/gin"
)

func MessageCreateRoute(c *gin.Context) {
	// Validating input
	req := &lib.MessageCreateSchema{}
	rs := lib.APIResponse{}
	if err := c.BindJSON(req); err != nil {
		er := true
		er_msg := "Invalid input"

		rs.Err = &er
		rs.Message = &er_msg
		c.JSON(http.StatusNotAcceptable, rs)
		return
	}

	room_id := c.Param("room_id")
	content := req.Content

	// TODO: Use snowflake ID and implement proper payload structuring
	payload := map[string]any{
		"content": content,
		"room_id": room_id,
	}

	rs.Data = &payload

	// Dispatch new message to NATS
	// TODO: Dispatch message to Kafka logs for batch processing
	lib.SendMsgNATS(payload)

	c.JSON(http.StatusCreated, rs)
}
