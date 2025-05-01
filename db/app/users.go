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

type UserDB struct {
	AppDB *db.AppDB
}

// Will be stubbed during testing
type IUserDB interface {
	CreateUser(ctx context.Context, id int64, username, displayName, email string, password []byte) error
	GetUserByAnyField(ctx context.Context, fields map[lib.UserColumn]any) (*lib.UserModel, error)
}

func GetUserDB() *UserDB {
	user_db := &UserDB{AppDB: db.GetAppDB()}
	return user_db
}

func (u *UserDB) CreateUser(ctx context.Context, id int64, username, displayName, email string, password []byte) error {
	/* Service function to create a user */

	// TODO: Separate DB functionality like insert, delete, etc. into separate functions
	tx, err := u.AppDB.DBPool.BeginTx(ctx, pgx.TxOptions{})
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

func (u *UserDB) GetUserByAnyField(ctx context.Context, fields map[lib.UserColumn]any) (*lib.UserModel, error) {
	/* This function queries user based on OR query for multiple fields */
	// TODO:

	// Separate the fields and respective values into properly sequenced slices for later use in query generation
	field_arr := []lib.UserColumn{}
	value_arr := []any{}
	for k, v := range fields {
		field_arr = append(field_arr, k)
		value_arr = append(value_arr, v)
	}

	// Generate the "OR" conditions using given fields
	or_fields := lib.GenerateDBOrFields(field_arr)

	// Generate the SELECT columns
	cols := []lib.UserColumn{lib.ColUserID, lib.ColUserUsername, lib.ColUserDisplayName, lib.ColUserCreatedAt, lib.ColUserEmail, lib.ColUserPassword}
	selected_columns := lib.GenerateDBQueryFields(cols)

	query := fmt.Sprintf("SELECT %s FROM %s INNER JOIN %s ON %s.%s=%s.%s WHERE ",
		selected_columns,
		lib.TableUsers,
		lib.TableAuth,
		lib.ColUserID.Table,
		lib.ColUserID.Column,
		lib.ColUserEmail.Table,
		lib.ColUserID.Column)
	query = query + or_fields

	// The value_arr maintains same sequence of parameters as the columns, which is why we separated the map into two slices
	row := u.AppDB.DBPool.QueryRow(ctx, query, value_arr...)
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

func insertUser(tx pgx.Tx, ctx context.Context, id int64, username, display_name string, created_at int64) error {
	_, err := tx.Exec(ctx, "INSERT INTO users (id, username, display_name, created_at) VALUES ($1, $2, $3, $4)", id, username, display_name, created_at)
	return err
}

func insertUsersAuth(tx pgx.Tx, ctx context.Context, id int64, email, pw_hash string) error {
	_, err := tx.Exec(ctx, "INSERT INTO users_auth (id, email, pw_hash) VALUES ($1, $2, $3)", id, email, pw_hash)
	return err
}
