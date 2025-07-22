package db

type DBColumn struct {
	// This struct includes both Column and Table because both may need to be specified for JOIN queries such as `user.id=user_auth.id`
	Column string
	Table  string
}

type ConversationType string
type ConversationRole string

const (
	CONVERSATION_DM    ConversationType = "DM"
	CONVERSATION_GROUP ConversationType = "GROUP"

	CONVERSATION_ROLE_ADMIN  ConversationRole = "ADMIN"
	CONVERSATION_ROLE_MEMBER ConversationRole = "MEMBER"
)
