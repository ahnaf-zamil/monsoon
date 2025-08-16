package db

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"strings"

	"monsoon/util"
)

const sqlFilesDir = "./schema"

func createDBSchemas(pool IPgxPool, schemaDir string) {
	files, err := os.ReadDir(schemaDir)
	if err != nil {
		log.Println("unable to read directory: ", err)
		return
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) != ".sql" {
			continue
		}

		filePath := filepath.Join(schemaDir, file.Name())
		sqlContent, err := os.ReadFile(filePath)
		if err != nil {
			log.Println("Unable to read file:", file.Name(), err)
			return
		}

		// Split statements by semicolon
		statements := strings.Split(string(sqlContent), ";")

		for _, stmt := range statements {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" {
				continue
			}

			_, err := pool.Exec(context.Background(), stmt)
			if err != nil {
				log.Println("Error executing statement from file:", file.Name(), err)
				return
			}
		}

		log.Println("Generated schema for:", file.Name())
	}
}

func CreateDBSchemas(conf *util.Config) {
	log.Println("Implementing all table schemas")
	if err := CreateAppDBPool(conf.AppDBPostgresURL); err != nil {
		panic(err)
	}
	defer CloseConnection()
	pool := GetAppDB().DBPool
	createDBSchemas(pool, sqlFilesDir)
}
