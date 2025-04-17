package main

import (
	"flag"
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
)

var PORT string = "8080"

func main() {
	// Load dotenv and config
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	conf := lib.LoadConfig()

	schemaFlagPtr := flag.Bool("generate_schema", false, "Generate DB schema from SQL files")
	flag.Parse()

	// Generate DB schema if the binary is run with this flag

	if *schemaFlagPtr {
		db.CreateAppDBSchemas(conf)
		return
	}

	// Initialize snowflake ID generator
	lib.InitSnowflakeNode()

	// Connect to database

	appDB := db.CreateAppDB()
	if err := appDB.CreateConnectionPool(conf.AppDBPostgresURL); err != nil {
		panic(err)
	}
	defer appDB.CloseConnection()

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
	log.Println("Server started on port", PORT)

	err := http.ListenAndServe("0.0.0.0:"+PORT, r)
	if err != nil {
		log.Println("Error starting server:", err)
	}
}
