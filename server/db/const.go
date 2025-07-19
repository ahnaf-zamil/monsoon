package db

/* Holds DB constants such as tables and columns for dynamic query generation

May refactor/rename this and place it somewhere else later, for now it stays here */

type UserColumn struct {
	// This struct includes both Column and Table because both may need to be specified for JOIN queries such as `user.id=user_auth.id`
	Column string
	Table  string
}

const (
	TableUsers        string = "users"
	TableAuth         string = "users_auth"
	TableUsersSession string = "users_session"
	TableUsersKey     string = "users_key"
)

var (
	// UserColumn -> {Column Name, Table Name}
	ColUserID          UserColumn = UserColumn{"id", TableUsers}
	ColUserUsername    UserColumn = UserColumn{"username", TableUsers}
	ColUserDisplayName UserColumn = UserColumn{"display_name", TableUsers}
	ColUserCreatedAt   UserColumn = UserColumn{"created_at", TableUsers}
	ColUserUpdatedAt   UserColumn = UserColumn{"updated_at", TableUsers}
	ColUserEmail       UserColumn = UserColumn{"email", TableAuth}
	ColUserPassword    UserColumn = UserColumn{"pw_hash", TableAuth}

	ColSessionID           UserColumn = UserColumn{"session_id", TableUsersSession}
	ColSessionUserID       UserColumn = UserColumn{"user_id", TableUsersSession}
	ColSessionRefreshToken UserColumn = UserColumn{"refresh_token", TableUsersSession}
	ColSessionCreatedAt    UserColumn = UserColumn{"created_at", TableUsersSession}

	ColKeyUserID UserColumn = UserColumn{"session_id", TableUsersKey}
	ColKeyEncKey UserColumn = UserColumn{"pub_enc_key", TableUsersKey}
	ColKeySigKey UserColumn = UserColumn{"pub_sig_key", TableUsersKey}
)
