package main

import (
	"log"

	"monsoon/db"
	"monsoon/util"

	"github.com/joho/godotenv"
)

/* A simple DB migration script */

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	conf := util.LoadConfig()
	db.CreateAppDBSchemas(conf)
	db.CreateMsgDBSchemas(conf)
}
