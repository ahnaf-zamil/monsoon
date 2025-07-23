package db

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"monsoon/util"
)

const appSchemaDir = "./schema/app"
const msgSchemaDir = "./schema/message"

func createDBSchemas(pool IPgxPool, schemaDir string) {
	files, err := os.ReadDir(schemaDir) // Path of the app DB sql files
	if err != nil {
		log.Println("unable to read directory: ", err)
		return
	}
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".sql" {
			filePath := filepath.Join(schemaDir, file.Name())

			sqlContent, err := os.ReadFile(filePath)
			if err != nil {
				log.Println("Unable to read file:", file.Name(), err)
				return
			}
			_, err = pool.Exec(context.Background(), string(sqlContent))
			if err != nil {
				log.Println("Error executing SQL from file:", file.Name(), err)
				return
			}
			log.Println("Generated schema for:", file.Name())
		}
	}
}

func CreateAppDBSchemas(conf *util.Config) {
	log.Println("Starting schema generation for App DB")
	if err := CreateAppDBPool(conf.AppDBPostgresURL); err != nil {
		panic(err)
	}
	defer CloseConnection()
	pool := GetAppDB().DBPool
	createDBSchemas(pool, appSchemaDir)
}

func CreateMsgDBSchemas(conf *util.Config) {
	log.Println("Starting schema generation for Message DB")
	if err := CreateMsgDBPool(conf.MessageDBPostgresURL); err != nil {
		panic(err)
	}
	defer CloseConnection()
	pool := GetMsgDB().DBPool
	createDBSchemas(pool, msgSchemaDir)
}
