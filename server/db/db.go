package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var userDBPool *pgxpool.Pool
var msgDBPool *pgxpool.Pool

// Stub during testing
type IPgxPool interface {
	BeginTx(context.Context, pgx.TxOptions) (pgx.Tx, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...any) pgx.Row
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
	Close()
}

type AppDB struct {
	DBPool IPgxPool
}

type MsgDB struct {
	DBPool IPgxPool
}

func CreateConnectionPool(dbURL string) (*pgxpool.Pool, error) {
	conn, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Printf("Unable to create connection pool: %v\n", err)
		return nil, err
	}
	err = conn.Ping(context.Background())
	if err != nil {
		log.Printf("DB Ping error: %v\n", err)
		return nil, err
	}
	return conn, nil
}

func CreateAppDBPool(dbURL string) error {
	pool, err := CreateConnectionPool(dbURL)
	if err != nil {
		return err
	}
	userDBPool = pool
	log.Println("Created App database connection pool")
	return nil
}

func CreateMsgDBPool(dbURL string) error {
	pool, err := CreateConnectionPool(dbURL)
	if err != nil {
		return err
	}
	msgDBPool = pool
	log.Println("Created Message database connection pool")
	return nil
}

func CloseConnection() {
	if userDBPool != nil {
		userDBPool.Close()
	}

	if msgDBPool != nil {
		msgDBPool.Close()
	}
}

func GetAppDB() *AppDB {
	return &AppDB{DBPool: userDBPool}
}

func GetMsgDB() *MsgDB {
	return &MsgDB{DBPool: msgDBPool}
}
