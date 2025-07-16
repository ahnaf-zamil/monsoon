package ws

import (
	"log"
	"net/http"

	"monsoon/lib"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Handles socket client connecting to the server
func RegisterSocketClient(wsConn *websocket.Conn, userId string) (*Socket, error) {
	s := &Socket{ID: lib.GenerateSocketID(), Rooms: make(map[string]bool), UserID: userId, WsConn: wsConn}
	AddSocketToList(s)

	// Adding socket to rooms
	// TODO: Connect database and fetch room list to join
	AddSocketToRoom(s, "amogus")

	// Return list of rooms to the client
	roomsList := []string{}
	for k := range s.Rooms {
		roomsList = append(roomsList, k)
	}
	// Handle rooms list send error
	// Client will keep all these rooms list in memory
	if err := s.WsConn.WriteJSON(map[string]any{"rooms": roomsList}); err != nil {
		defer s.WsConn.Close()
		return s, err
	}

	return s, nil
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

func WsHandler(c *gin.Context) {
	// Upgrading connection to websocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Error upgrading:", err)
		return
	}

	defer conn.Close()

	log.Println("Client connected from:", c.ClientIP())

	// Take authentication token. For now we consider token to be user ID
	token := c.Query("token")

	if token == "" {
		if err := conn.WriteJSON(map[string]any{"error": "Unauthorized: Missing token"}); err != nil {
			log.Println("Error writing msg:", err)
		}

		log.Printf("Disconnect: %s connected with bad token", c.ClientIP())
		return
	}

	// TODO: Fetch user ID by using token for Socket struct init
	// For now the token is being used as user ID
	userId := token

	// Add client to socket state
	client_s, err := RegisterSocketClient(conn, userId)
	// Only continue working with client if registration/initialization is successful
	if err == nil {
		// PrintSocketList()
		// PrettyPrintSyncMap(GetRoomState())
		for {
			msg_data := map[string]any{}
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

			log.Println("Recv:", msg_data)
		}
	}
	// Handle disconnection to cleanup socket and room states
	HandleSocketDisconnect(client_s)
}
