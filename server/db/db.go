package db

import (
	"context"
	"log"
	"monsoon/util"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var userDBPool *pgxpool.Pool

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

type QueryLogger struct{}

func (l QueryLogger) TraceQueryStart(
	ctx context.Context,
	_ *pgx.Conn,
	data pgx.TraceQueryStartData,
) context.Context {
	log.Printf("[QUERY START] SQL: %s, Args: %v", data.SQL, data.Args)
	return ctx
}

func (l QueryLogger) TraceQueryEnd(
	ctx context.Context,
	_ *pgx.Conn,
	data pgx.TraceQueryEndData) {
	if data.Err != nil {
		log.Printf("[QUERY END] ERROR: %v", data.Err)
	} else {
		log.Printf("[QUERY END] Success")
	}
}

func CreateConnectionPool(dbURL string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		log.Fatal(err)
	}

	if util.IsDevEnv() {
		config.ConnConfig.Tracer = QueryLogger{}
	}

	conn, err := pgxpool.NewWithConfig(context.Background(), config)
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
	log.Println("Created database connection pool")
	return nil
}

func CloseConnection() {
	if userDBPool != nil {
		userDBPool.Close()
	}
}

func GetAppDB() *AppDB {
	return &AppDB{DBPool: userDBPool}
}
