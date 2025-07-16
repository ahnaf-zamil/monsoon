package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var pool *pgxpool.Pool

// Stub during testing
type IPgxPool interface {
	BeginTx(context.Context, pgx.TxOptions) (pgx.Tx, error)
	QueryRow(context.Context, string, ...any) pgx.Row
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
	Close()
}

type AppDB struct {
	DBPool IPgxPool
}

func CreateConnectionPool(dbURL string) error {
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
	pool = conn
	return nil
}

func CloseConnection() {
	pool.Close()
}

func GetAppDB() *AppDB {
	return &AppDB{DBPool: pool}
}
