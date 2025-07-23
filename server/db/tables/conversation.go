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

	ColConvoPartConversationID db.DBColumn = db.DBColumn{Column: "conversation_id", Table: TableConversationParticipants}
	ColConvoPartUserID         db.DBColumn = db.DBColumn{Column: "user_id", Table: TableConversationParticipants}
	ColConvoPartJoinedAt       db.DBColumn = db.DBColumn{Column: "joined_at", Table: TableConversationParticipants}
	ColConvoPartRole           db.DBColumn = db.DBColumn{Column: "conversatroleion_id", Table: TableConversationParticipants}
)
