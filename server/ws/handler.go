package ws

import (
	"log"
	"net/http"
	"time"

	"monsoon/db/app"
	"monsoon/lib"
	"monsoon/util"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func GetWebSocketHandler() IWebSocketHandler {
	return &WebSocketHandler{UserDB: app.GetUserDB(), TokenHelper: lib.GetJWTTokenHelper()}
}

func (w *WebSocketHandler) HandleSocketEvent(socket *Socket, msg map[string]any) {
	if opcodeRaw, ok := msg["opcode"]; ok {
		if opcodeStr, ok := opcodeRaw.(string); ok {
			switch EventOpCode(opcodeStr) {
			case OpHeartbeat:
				w.HandleClientHeartbeat(socket)
			}
		}
	}
}

// Handles socket client connecting to the server
func (w *WebSocketHandler) RegisterSocketClient(wsConn *websocket.Conn, userID string) (*Socket, error) {
	s := &Socket{ID: util.GenerateSocketID(), Rooms: make(map[string]bool), UserID: userID, WsConn: wsConn, LastHeartbeat: time.Now()}
	mu.RLock()
	AddSocketToList(s)

	err := w.DispatchEvent(s, OpHeartbeatInit, gin.H{"interval": HEARTBEAT_INTERVAL.Milliseconds(), "timeout": HEARTBEAT_TIMEOUT.Milliseconds()})
	if err != nil {
		RemoveSocketFromList(s)
		return nil, err
	}
	mu.RUnlock()

	// Adding socket to rooms
	// TODO: Connect database and fetch room list to join

	// Add user to own DM room to receive DMs
	AddSocketToRoom(s, "dm:"+userID)

	// Return list of rooms to the client
	roomsList := []string{}
	for k := range s.Rooms {
		roomsList = append(roomsList, k)
	}
	// Handle rooms list send error
	// Client will keep all these rooms list in memory
	err = w.DispatchEvent(s, OpRoomSync, gin.H{"rooms": roomsList})
	return s, err
}

// Handles any socket disconnection event
func HandleSocketDisconnect(client_s *Socket) {
	// Removes socket from the sock list state
	log.Printf("Disconnected client: %s (%s)\n", client_s.ID, client_s.WsConn.RemoteAddr())
	RemoveSocketFromList(client_s)

	// Remove socket from all rooms
	for k := range client_s.Rooms {
		RemoveSocketFromRoom(client_s, k)
	}
}

func (w *WebSocketHandler) ConnectionHandler(c *gin.Context) {
	// Upgrading connection to websocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Error upgrading:", err)
		return
	}

	defer conn.Close()

	log.Println("Client connected from:", c.ClientIP())
	token := c.Query("token")

	if token == "" {
		if err := conn.WriteJSON(map[string]any{"error": "Unauthorized: Missing token"}); err != nil {
			log.Println("Error writing msg:", err)
		}

		log.Printf("Disconnect: %s connected with bad token", c.ClientIP())
		return
	}

	userID, err := w.TokenHelper.VerifyToken(token)
	if err != nil {
		if err := conn.WriteJSON(map[string]any{"error": "Unauthorized: Missing token"}); err != nil {
			log.Println("Error writing msg:", err)
		}

		log.Printf("Disconnect: %s connected with bad token", c.ClientIP())
		return
	}

	user, _ := w.UserDB.GetUserByID(c.Request.Context(), userID)
	if user == nil {
		if err := conn.WriteJSON(map[string]any{"error": "Unauthorized"}); err != nil {
			log.Println("Error writing msg:", err)
		}
	}

	// Add client to socket state
	client_s, err := w.RegisterSocketClient(conn, userID)
	// Only continue working with client if registration/initialization is successful
	if err == nil {
		// util.PrettyPrintSyncMap(GetRoomState())
		for {
			msg_data := gin.H{}
			err := client_s.WsConn.ReadJSON(&msg_data)
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
					log.Println("Connection closed unexpectedly:", err)
				} else if err.Error() == "websocket: read on closed connection" {
					log.Println("Connection was closed gracefully by client")
				} else {
					log.Println("Error reading JSON:", err)
				}
				break // Exit the loop if the connection has an error or is closed
			}

			w.HandleSocketEvent(client_s, msg_data)
		}
	}
	// Handle disconnection to cleanup socket and room states
	HandleSocketDisconnect(client_s)
}
