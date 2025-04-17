package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var appDB *AppDB

type AppDB struct {
	DBPool *pgxpool.Pool
}

func (app_db *AppDB) CreateConnectionPool(dbURL string) error {
	conn, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Printf("Unable to create connection pool: %v\n", err)
		return err
	}
	err = conn.Ping(context.Background())
	if err != nil {
		log.Printf("DB Ping error: %v\n", err)
		return err
	}
	log.Println("Connected to database")

	app_db.DBPool = conn
	return nil
}

func (app_db *AppDB) CloseConnection() {
	app_db.DBPool.Close()
}

func GetAppDB() *AppDB {
	return appDB
}

func CreateAppDB() *AppDB {
	appDB = &AppDB{}
	return appDB
}
