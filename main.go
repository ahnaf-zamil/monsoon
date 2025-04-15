package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"ws_realtime_app/controller"
	"ws_realtime_app/lib"
	"ws_realtime_app/ws"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Handles socket client connecting to the server
func RegisterSocketClient(wsConn *websocket.Conn, userId string) (*ws.Socket, error) {
	s := &ws.Socket{ID: lib.GenerateSocketID(), Rooms: make(map[string]bool), UserID: userId, WsConn: wsConn}
	ws.AddSocketToList(s)

	// Adding socket to rooms
	// TODO: Connect database and fetch room list to join
	ws.AddSocketToRoom(s, "amogus")

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
func HandleSocketDisconnect(client_s *ws.Socket) {
	// Removes socket from the sock list state
	log.Printf("Disconnected client: %s (%s)\n", client_s.ID, client_s.WsConn.RemoteAddr())
	ws.RemoveSocketFromList(client_s)

	// Remove socket from all rooms
	for k := range client_s.Rooms {
		ws.RemoveSocketFromRoom(client_s, k)
	}
}

func WsHandler(c *gin.Context) {
	// Upgrading connection to websocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("Error upgrading:", err)
		return
	}

	defer conn.Close()

	log.Println("Client connected from:", c.ClientIP())

	// Take authentication token. For now we consider token to be user ID
	token := c.Query("token")

	if token == "" {
		if err := conn.WriteJSON(map[string]any{"error": "Unauthorized: Missing token"}); err != nil {
			fmt.Println("Error writing msg:", err)
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

func main() {
	// Load dotenv and config
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	conf := lib.LoadConfig()

	// Initialize snowflake ID generator
	lib.InitSnowflakeNode()

	// Initialize NATS connection and drain connection upon return
	nc := lib.InitNATS(conf.NATSUrl)
	defer nc.Drain()

	// Initialize Gin with CORS, and register controller routes
	r := gin.Default()
	r.GET("/ws", WsHandler)
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	}))
	controller.InitControllers(r)

	// Here we go
	fmt.Println("Websocket server started")

	err := http.ListenAndServe("0.0.0.0:8080", r)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
