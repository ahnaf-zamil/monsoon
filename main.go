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
	"github.com/joho/godotenv"
)

func main() {
	// Load dotenv and config
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	conf := lib.LoadConfig()

	// Initialize snowflake ID generator
	lib.InitSnowflakeNode()

	// Initialize NATS connection and drain connection upon return
	nc := ws.InitNATS(conf.NATSUrl)
	defer nc.Drain()

	// Initialize Gin with CORS, and register controller routes
	r := gin.Default()
	r.GET("/ws", ws.WsHandler)
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
