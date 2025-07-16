package db

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"monsoon/lib"
)

var appSchemaDir = "./schema/app"

func CreateAppDBSchemas(conf *lib.Config) {
	log.Println("Starting schema generation for App DB")

	if err := CreateConnectionPool(conf.AppDBPostgresURL); err != nil {
		panic(err)
	}
	defer CloseConnection()

	pool := GetAppDB().DBPool

	files, err := os.ReadDir(appSchemaDir) // Path of the app DB sql files
	if err != nil {
		log.Println("unable to read directory: ", err)
		return
	}
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".sql" {
			filePath := filepath.Join(appSchemaDir, file.Name())

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
