package main

import (
	"log"
	"net/http"
	"time"

	"monsoon/controller"
	"monsoon/db"
	"monsoon/lib"
	"monsoon/util"
	"monsoon/ws"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	swaggerfiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware

	_ "monsoon/docs" // Loading swagger docs
)

//	@title			Monsoon API
//	@version		0.0.1
//	@description REST API and WebSocket server for Monsoon

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

	conf := util.LoadConfig()

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
	if !(conf.IsDev) {
		gin.SetMode(gin.ReleaseMode)
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "https://monsoon.ahnafzamil.com"},
		AllowHeaders:     []string{"content-type", "authorization"},
		AllowMethods:     []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
		AllowCredentials: true,
		MaxAge:           24 * 7 * time.Hour,
	}))

	controller.InitControllers(r)

	wsHandler := ws.GetWebSocketHandler()
	r.GET("/ws", wsHandler.ConnectionHandler)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// Here we go
	log.Println("Server started on port", conf.Port)

	wsHandler.StartHeartbeat()
	err = http.ListenAndServe("0.0.0.0:"+conf.Port, r)
	if err != nil {
		log.Println("Error starting server:", err)
	}
}
