package ws

import (
	"monsoon/db/app"
	"monsoon/lib"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type IWebSocketHandler interface {
	ConnectionHandler(c *gin.Context)
	RegisterSocketClient(wsConn *websocket.Conn, userID string) (*Socket, error)
	StartHeartbeat()
	DispatchEvent(socket *Socket, opCode EventOpCode, data any) error
}

type WebSocketHandler struct {
	UserDB      app.IUserDB
	TokenHelper lib.IJWTTokenHelper
}

type EventOpCode string

type WebSocketEvent struct {
	OpCode EventOpCode `json:"opcode"`
	Data   any         `json:"data,omitempty"`
}
