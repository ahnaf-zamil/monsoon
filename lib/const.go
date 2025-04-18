package lib

/* Holds DB constants such as tables and columns for dynamic query generation

May refactor/rename this and place it somewhere else later, for now it stays here */

type UserColumn struct {
	// This struct includes both Column and Table because both may need to be specified for JOIN queries such as `user.id=user_auth.id`
	Column string
	Table  string
}

const (
	TableUsers string = "users"
	TableAuth  string = "users_auth"
)

var (
	// UserColum -> {Column Name, Table Name}
	ColUserID          UserColumn = UserColumn{"id", TableUsers}
	ColUserUsername    UserColumn = UserColumn{"username", TableUsers}
	ColUserDisplayName UserColumn = UserColumn{"display_name", TableUsers}
	ColUserCreatedAt   UserColumn = UserColumn{"created_at", TableUsers}
	ColUserEmail       UserColumn = UserColumn{"email", TableAuth}
	ColUserPassword    UserColumn = UserColumn{"pw_hash", TableAuth}
)
