package tables

import "monsoon/db"

const (
	TableConversations            = "conversations"
	TableDirectConversations      = "direct_conversations"
	TableConversationParticipants = "conversation_participants"
)

var (
	ColDMID        db.DBColumn = db.DBColumn{Column: "conversation_id", Table: TableDirectConversations}
	ColDMUser1     db.DBColumn = db.DBColumn{Column: "user1", Table: TableDirectConversations}
	ColDMUser2     db.DBColumn = db.DBColumn{Column: "user2", Table: TableDirectConversations}
	ColDMCreatedAt db.DBColumn = db.DBColumn{Column: "created_at", Table: TableDirectConversations}
)
