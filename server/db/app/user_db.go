package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"monsoon/api"
	"monsoon/db"
	"time"

	"github.com/jackc/pgx/v5"
)

type UserDB struct {
	AppDB *db.AppDB
}

// Will be stubbed during testing
type IUserDB interface {
	CreateUser(ctx context.Context, id int64, username, displayName, email string, password []byte, keyEnc []byte, keySig []byte, pwSalt []byte, encSalt []byte, encSeed []byte, nonce []byte) error
	GetUserByAnyField(ctx context.Context, fields map[db.UserColumn]any) (*api.UserModel, error)
	UpdateUserTableById(ctx context.Context, id int64, table string, values map[db.UserColumn]string) error
	CreateUserSession(ctx context.Context, sessionID int64, userID int64, refreshToken string) error
	GetSessionByAnyField(ctx context.Context, fields map[db.UserColumn]any) (*api.UserSessionModel, error)
	GetUserAuthByID(ctx context.Context, userID string) (*api.UserAuthModel, error)
	// Util
	GetUserByID(c context.Context, userID string) (*api.UserModel, error)
}

func GetUserDB() *UserDB {
	user_db := &UserDB{AppDB: db.GetAppDB()}
	return user_db
}

func (u *UserDB) CreateUser(ctx context.Context, id int64, username, displayName, email string, password []byte, keyEnc []byte, keySig []byte, pwSalt []byte, encSalt []byte, encSeed []byte, nonce []byte) error {
	/* Service function to create a user
	 */

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
	err = insertUserAuth(tx, ctx, id, email, password, pwSalt, encSalt, encSeed, nonce)
	if err != nil {
		return err
	}

	// Insert user's signature and encryption key in key table
	err = insertUserKey(tx, ctx, id, keyEnc, keySig)
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
	cols := []db.UserColumn{db.ColUserID, db.ColUserUsername, db.ColUserDisplayName, db.ColUserCreatedAt, db.ColUserEmail, db.ColUserPassword}
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
	err := row.Scan(&user.ID, &user.Username, &user.DisplayName, &user.CreatedAt, &user.Email, &user.PasswordHash)
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

func (u *UserDB) CreateUserSession(ctx context.Context, sessionID int64, userID int64, refreshToken string) error {
	/* Service function to create a user session */

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
	sessionCreated := time.Now().Unix()
	err = insertUserSession(tx, ctx, sessionID, userID, refreshToken, sessionCreated)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Println("TX commit error:", err)
	}
	return nil
}

func (u *UserDB) GetSessionByAnyField(ctx context.Context, fields map[db.UserColumn]any) (*api.UserSessionModel, error) {
	/* This function queries user session based on OR query for multiple fields */

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
	cols := []db.UserColumn{db.ColSessionID, db.ColSessionUserID, db.ColSessionRefreshToken, db.ColSessionCreatedAt}
	selected_columns := db.GenerateDBQueryFields(cols)

	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s",
		selected_columns,
		db.TableUsersSession,
		or_fields)

	// The value_arr maintains same sequence of parameters as the columns, which is why we separated the map into two slices
	row := u.AppDB.DBPool.QueryRow(ctx, query, value_arr...)

	var session api.UserSessionModel
	err := row.Scan(&session.SessionID, &session.UserID, &session.RefreshToken, &session.CreatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		// Return nil if no rows
		return nil, nil
	} else if err != nil {
		// If there's error, just return error
		return nil, err
	} else {
		// Else, just return user
		return &session, nil
	}
}

func (u *UserDB) GetUserAuthByID(ctx context.Context, userID string) (*api.UserAuthModel, error) {
	// TODO: implement later
	cols := []db.UserColumn{db.ColAuthUserID, db.ColAuthPasswordHash, db.ColAuthPasswordSalt, db.ColAuthEncryptionSalt, db.ColAuthEncryptedSeed, db.ColAuthNonce}
	selected_columns := db.GenerateDBQueryFields(cols)

	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s = $1",
		selected_columns,
		db.TableAuth, db.ColAuthUserID.Column)
	// The value_arr maintains same sequence of parameters as the columns, which is why we separated the map into two slices
	row := u.AppDB.DBPool.QueryRow(ctx, query, userID)
	var authModel api.UserAuthModel
	err := row.Scan(&authModel.UserID, &authModel.PasswordHash, &authModel.PasswordSalt, &authModel.EncryptionSalt, &authModel.EncryptedSeed, &authModel.Nonce)

	if errors.Is(err, pgx.ErrNoRows) {
		// Return nil if no rows
		return nil, nil
	} else if err != nil {
		// If there's error, just return error
		return nil, err
	} else {
		// Else, just return user
		return &authModel, nil
	}
}

func insertUser(tx pgx.Tx, ctx context.Context, id int64, username, display_name string, created_at int64, updated_at int64) error {
	query := fmt.Sprintf("INSERT INTO %s (id, username, display_name, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)", db.TableUsers)
	_, err := tx.Exec(ctx, query, id, username, display_name, created_at, updated_at)
	return err
}

func insertUserAuth(tx pgx.Tx, ctx context.Context, id int64, email string, pw_hash []byte, pwSalt []byte, encSalt []byte, encSeed []byte, nonce []byte) error {
	query := fmt.Sprintf("INSERT INTO %s (id, email, pw_hash, pw_salt, enc_salt, key_seed_cipher, nonce) VALUES ($1, $2, $3, $4, $5, $6, $7)", db.TableAuth)
	_, err := tx.Exec(ctx, query, id, email, pw_hash, pwSalt, encSalt, encSeed, nonce)
	return err
}

func insertUserSession(tx pgx.Tx, ctx context.Context, sessionID int64, userID int64, refreshToken string, sessionCreated int64) error {
	query := fmt.Sprintf("INSERT INTO %s (session_id, user_id, refresh_token, created_at) VALUES ($1, $2, $3, $4)", db.TableUsersSession)
	_, err := tx.Exec(ctx, query, sessionID, userID, refreshToken, sessionCreated)
	return err
}

func insertUserKey(tx pgx.Tx, ctx context.Context, userID int64, keyEnc []byte, keySig []byte) error {
	query := fmt.Sprintf("INSERT INTO %s (user_id, pub_sig_key, pub_enc_key) VALUES ($1, $2, $3)", db.TableUsersKey)
	_, err := tx.Exec(ctx, query, userID, keySig, keyEnc)
	return err
}
