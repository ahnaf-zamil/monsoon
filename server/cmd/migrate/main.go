package main

import (
	"log"
	"monsoon/db"
	"monsoon/lib"

	"github.com/joho/godotenv"
)

/* A simple DB migration script */

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	conf := lib.LoadConfig()
	db.CreateAppDBSchemas(conf)
}
