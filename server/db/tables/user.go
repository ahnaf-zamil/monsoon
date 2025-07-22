package tables

import "monsoon/db"

/* Holds DB constants such as tables and columns for dynamic query generation

May refactor/rename this and place it somewhere else later, for now it stays here */

const (
	TableUsers        string = "users"
	TableAuth         string = "users_auth"
	TableUsersSession string = "users_session"
	TableUsersKey     string = "users_key"
)

var (
	// db.DBColumn -> {Column Name, Table Name}
	ColUserID          db.DBColumn = db.DBColumn{Column: "id", Table: TableUsers}
	ColUserUsername    db.DBColumn = db.DBColumn{Column: "username", Table: TableUsers}
	ColUserDisplayName db.DBColumn = db.DBColumn{Column: "display_name", Table: TableUsers}
	ColUserCreatedAt   db.DBColumn = db.DBColumn{Column: "created_at", Table: TableUsers}
	ColUserUpdatedAt   db.DBColumn = db.DBColumn{Column: "updated_at", Table: TableUsers}
	ColUserEmail       db.DBColumn = db.DBColumn{Column: "email", Table: TableAuth}
	ColUserPassword    db.DBColumn = db.DBColumn{Column: "pw_hash", Table: TableAuth}

	ColSessionID           db.DBColumn = db.DBColumn{Column: "session_id", Table: TableUsersSession}
	ColSessionUserID       db.DBColumn = db.DBColumn{Column: "user_id", Table: TableUsersSession}
	ColSessionRefreshToken db.DBColumn = db.DBColumn{Column: "refresh_token", Table: TableUsersSession}
	ColSessionCreatedAt    db.DBColumn = db.DBColumn{Column: "created_at", Table: TableUsersSession}

	ColKeyUserID db.DBColumn = db.DBColumn{Column: "user_id", Table: TableUsersKey}
	ColKeyEncKey db.DBColumn = db.DBColumn{Column: "pub_enc_key", Table: TableUsersKey}
	ColKeySigKey db.DBColumn = db.DBColumn{Column: "pub_sig_key", Table: TableUsersKey}

	ColAuthUserID         db.DBColumn = db.DBColumn{Column: "id", Table: TableAuth}
	ColAuthPasswordHash   db.DBColumn = db.DBColumn{Column: "pw_hash", Table: TableAuth}
	ColAuthPasswordSalt   db.DBColumn = db.DBColumn{Column: "pw_salt", Table: TableAuth}
	ColAuthEncryptionSalt db.DBColumn = db.DBColumn{Column: "enc_salt", Table: TableAuth}
	ColAuthEncryptedSeed  db.DBColumn = db.DBColumn{Column: "key_seed_cipher", Table: TableAuth}
	ColAuthNonce          db.DBColumn = db.DBColumn{Column: "nonce", Table: TableAuth}
)
