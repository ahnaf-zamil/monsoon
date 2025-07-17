package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"monsoon/api"
	"monsoon/db"

	"github.com/jackc/pgx/v5"
)

type UserDB struct {
	AppDB *db.AppDB
}

// Will be stubbed during testing
type IUserDB interface {
	CreateUser(ctx context.Context, id int64, username, displayName, email string, password []byte, refreshToken string) error
	GetUserByAnyField(ctx context.Context, fields map[db.UserColumn]any) (*api.UserModel, error)
	UpdateUserTableById(ctx context.Context, id int64, table string, values map[db.UserColumn]string) error
}

func GetUserDB() *UserDB {
	user_db := &UserDB{AppDB: db.GetAppDB()}
	return user_db
}

func (u *UserDB) CreateUser(ctx context.Context, id int64, username, displayName, email string, password []byte, refreshToken string) error {
	/* Service function to create a user */

	// TODO: Separate DB functionality like insert, delete, etc. into separate functions
	tx, err := u.AppDB.DBPool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != pgx.ErrTxClosed {
			log.Println("TX rollback error:", err)
		}
	}()

	// Insert user data in users table
	createdAt := time.Now().Unix()
	updatedAt := createdAt
	err = insertUser(tx, ctx, id, username, displayName, createdAt, updatedAt)
	if err != nil {
		return err
	}

	// Insert user auth data into the auth table
	err = insertUsersAuth(tx, ctx, id, email, string(password), refreshToken)
	if err != nil {
		return err
	}
	err = tx.Commit(ctx)
	if err != nil {
		log.Println("TX commit error:", err)
	}

	return nil
}

func (u *UserDB) GetUserByAnyField(ctx context.Context, fields map[db.UserColumn]any) (*api.UserModel, error) {
	/* This function queries user based on OR query for multiple fields */
	// TODO:

	// Separate the fields and respective values into properly sequenced slices for later use in query generation
	field_arr := []db.UserColumn{}
	value_arr := []any{}
	for k, v := range fields {
		field_arr = append(field_arr, k)
		value_arr = append(value_arr, v)
	}

	// Generate the "OR" conditions using given fields
	or_fields := db.GenerateDBOrFields(field_arr)

	// Generate the SELECT columns
	cols := []db.UserColumn{db.ColUserID, db.ColUserUsername, db.ColUserDisplayName, db.ColUserCreatedAt, db.ColUserEmail, db.ColUserPassword, db.ColUserRefreshToken}
	selected_columns := db.GenerateDBQueryFields(cols)

	query := fmt.Sprintf("SELECT %s FROM %s INNER JOIN %s ON %s.%s=%s.%s WHERE ",
		selected_columns,
		db.TableUsers,
		db.TableAuth,
		db.ColUserID.Table,
		db.ColUserID.Column,
		db.ColUserEmail.Table,
		db.ColUserID.Column)
	query = query + or_fields
	// The value_arr maintains same sequence of parameters as the columns, which is why we separated the map into two slices
	row := u.AppDB.DBPool.QueryRow(ctx, query, value_arr...)
	var user api.UserModel
	err := row.Scan(&user.ID, &user.Username, &user.DisplayName, &user.CreatedAt, &user.Email, &user.Password, &user.RefreshToken)
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

func (u *UserDB) UpdateUserTableById(ctx context.Context, id int64, table string, values map[db.UserColumn]string) error {
	tx, err := u.AppDB.DBPool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != pgx.ErrTxClosed {
			log.Println("TX rollback error:", err)
		}
	}()

	query_format := "UPDATE %s SET %s WHERE %s.id = $%d"
	query := fmt.Sprintf(query_format, table, db.GenerateDBUpdateFields(values), table, len(values)+1)
	value_arr := []any{}
	for _, v := range values {
		value_arr = append(value_arr, v)
	}
	_, err = tx.Exec(ctx, query, append(value_arr, id)...)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Println("TX commit error:", err)
	}

	return err
}

func insertUser(tx pgx.Tx, ctx context.Context, id int64, username, display_name string, created_at int64, updated_at int64) error {
	_, err := tx.Exec(ctx, "INSERT INTO users (id, username, display_name, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)", id, username, display_name, created_at, updated_at)
	return err
}

func insertUsersAuth(tx pgx.Tx, ctx context.Context, id int64, email, pw_hash string, refresh_token string) error {
	_, err := tx.Exec(ctx, "INSERT INTO users_auth (id, email, pw_hash, refresh_token) VALUES ($1, $2, $3, $4)", id, email, pw_hash, refresh_token)
	return err
}
