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
//  @license.name	AGPL-3.0
// 	@license.url	https://www.gnu.org/licenses/agpl-3.0.en.html
// 	@BasePath /api

// @produce json
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Load dotenv and config
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	var conf *util.Config
	if util.IsDevEnv() {
		// For dev, use .env config
		conf = util.LoadDotenvConfig()
	} else {
		// TODO: Use a shared configuration database for PROD e.g Hashicorp Vault. For now default to .env
		conf = util.LoadDotenvConfig()
	}

	// Initialize snowflake ID generator
	lib.InitSnowflakeNode()

	// Connect to app database
	if err := db.CreateAppDBPool(conf.AppDBPostgresURL); err != nil {
		panic(err)
	}
	defer db.CloseConnection()

	// Initialize Gin with CORS, and register controller routes
	r := gin.Default()
	if !(util.IsDevEnv()) {
		gin.SetMode(gin.ReleaseMode)
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     conf.AllowedOrigins,
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

	// Initialize NATS connection and drain connection upon return
	n := &ws.NATSPublisher{W: wsHandler}
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

	wsHandler.StartHeartbeat()
	err = http.ListenAndServe("0.0.0.0:"+conf.Port, r)
	if err != nil {
		log.Println("Error starting server:", err)
	}
}
