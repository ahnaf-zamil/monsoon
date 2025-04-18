package app

import (
	"context"
	"errors"
	"fmt"
	"time"
	"ws_realtime_app/db"
	"ws_realtime_app/lib"

	"github.com/jackc/pgx/v5"
)

type _UserDB struct{}

var Users _UserDB

func (u *_UserDB) CreateUser(ctx context.Context, id int64, username, displayName, email string, password []byte) error {
	/* Service function to create a user */

	// TODO: Separate DB functionality like insert, delete, etc. into separate functions
	app_db := db.GetAppDB()
	tx, err := app_db.DBPool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Insert user data in users table
	createdAt := time.Now().Unix()
	err = insertUser(tx, ctx, id, username, displayName, createdAt)
	if err != nil {
		return err
	}

	// Insert user auth data into the auth table
	err = insertUsersAuth(tx, ctx, id, email, string(password))
	if err != nil {
		return err
	}
	tx.Commit(ctx)
	return nil
}

func GetUserByField(ctx context.Context, column lib.UserColumn, value any) (*lib.UserModel, error) {
	/* This function queries user based on single column value */

	allowed := map[lib.UserColumn]bool{
		lib.ColUserID:       true,
		lib.ColUserEmail:    true,
		lib.ColUserUsername: true,
	}

	if !allowed[column] {
		return nil, fmt.Errorf("invalid query column: %s", column)
	}

	app_db := db.GetAppDB()

	// Dynamically generate select columns list

	cols := []lib.UserColumn{lib.ColUserID, lib.ColUserUsername, lib.ColUserDisplayName, lib.ColUserCreatedAt, lib.ColUserEmail, lib.ColUserPassword}
	selected_columns := lib.GenerateDBQueryFields(cols)
	// Looks pump in the all the goodies
	query := fmt.Sprintf("SELECT %s FROM %s INNER JOIN %s ON %s.%s=%s.%s WHERE ",
		selected_columns,
		lib.TableUsers,
		lib.TableAuth,
		lib.ColUserID.Table,
		lib.ColUserID.Column,
		lib.ColUserEmail.Table,
		lib.ColUserID.Column)
	query = query + fmt.Sprintf("%s.%s = $1", column.Table, column.Column)

	// BEHOLD, THE POWER OF RAW SQL
	row := app_db.DBPool.QueryRow(ctx, query, value)

	var user lib.UserModel
	err := row.Scan(&user.ID, &user.Username, &user.DisplayName, &user.CreatedAt, &user.Email, &user.Password)

	if errors.Is(err, pgx.ErrNoRows) {
		// Return nil if no rows
		return nil, nil
	} else if err != nil {
		// If there's error, just return error
		return nil, err
	} else {
		// Else, just return user
		return &user, nil
	}
}

func GetUserByAnyField(ctx context.Context, fields map[lib.UserColumn]any) (*lib.UserModel, error) {
	/* This function queries user based on OR query for multiple fields */
	// TODO:
	return nil, nil
}

func insertUser(tx pgx.Tx, ctx context.Context, id int64, username, display_name string, created_at int64) error {
	_, err := tx.Exec(ctx, "INSERT INTO users (id, username, display_name, created_at) VALUES ($1, $2, $3, $4)", id, username, display_name, created_at)
	return err
}

func insertUsersAuth(tx pgx.Tx, ctx context.Context, id int64, email, pw_hash string) error {
	_, err := tx.Exec(ctx, "INSERT INTO users_auth (id, email, pw_hash) VALUES ($1, $2, $3)", id, email, pw_hash)
	return err
}
