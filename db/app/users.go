package app

import (
	"context"
	"log"
	"time"
	"ws_realtime_app/db"

	"github.com/jackc/pgx/v5"
)

type _UserDB struct{}

var Users _UserDB

func (u *_UserDB) CreateUser(ctx context.Context, id int64, username, displayName, email string, password []byte) error {
	// TODO: Separate DB functionality like insert, delete, etc. into separate functions
	app_db := db.GetAppDB()
	tx, err := app_db.DBPool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Insert user data in users table
	createdAt := time.Now().Unix()
	_, err = tx.Exec(ctx, "INSERT INTO users (id, username, display_name, created_at) VALUES ($1, $2, $3, $4)", id, username, displayName, createdAt)
	if err != nil {
		return err
	}

	// Insert user auth data into the auth table
	_, err = tx.Exec(ctx, "INSERT INTO users_auth (id, email, pw_hash) VALUES ($1, $2, $3)", id, email, string(password))
	if err != nil {
		return err
	}
	c := tx.Commit(ctx)
	log.Println(c)
	return nil
}
