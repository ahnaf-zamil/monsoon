package main

import (
	"log"
	"net/http"
	"time"
	"ws_realtime_app/controller"
	"ws_realtime_app/db"
	"ws_realtime_app/lib"
	"ws_realtime_app/ws"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	swaggerfiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware

	_ "ws_realtime_app/docs" // Loading swagger docs
)

//	@title			WS_RT_APP API
//	@version		0.0.1
//	@description REST API and WebSocket server

//	@contact.name	Author
//	@contact.url	https://ahnafzamil.com/contact
//	@contact.email	ahnaf@ahnafzamil.com
// 	@BasePath /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Load dotenv and config
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	conf := lib.LoadConfig()

	// Initialize snowflake ID generator
	lib.InitSnowflakeNode()

	// Connect to database

	if err := db.CreateConnectionPool(conf.AppDBPostgresURL); err != nil {
		panic(err)
	}
	defer db.CloseConnection()

	// Initialize NATS connection and drain connection upon return
	n := &ws.NATSPublisher{}
	nc, err := n.InitNATS(conf.NATSUrl)
	if err != nil {
		panic(err)
	}
	defer func() {
		err := nc.Drain()
		if err != nil {
			log.Println("NATS connection drain error:", err)
		}
	}()

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
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	controller.InitControllers(r)

	// Here we go
	log.Println("Server started on port", conf.Port)

	err = http.ListenAndServe("0.0.0.0:"+conf.Port, r)
	if err != nil {
		log.Println("Error starting server:", err)
	}
}
