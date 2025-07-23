package tables

import "monsoon/db"

const (
	TableMessages = "messages"
)

var (
	ColMessageID             db.DBColumn = db.DBColumn{Column: "id", Table: TableMessages}
	ColMessageConversationID db.DBColumn = db.DBColumn{Column: "conversation_id", Table: TableMessages}
	ColMessageAuthorID       db.DBColumn = db.DBColumn{Column: "author_id", Table: TableMessages}

	ColMessageContent db.DBColumn = db.DBColumn{Column: "content", Table: TableMessages}

	ColMessageCreatedAt db.DBColumn = db.DBColumn{Column: "created_at", Table: TableMessages}
	ColMessageEditedAt  db.DBColumn = db.DBColumn{Column: "edited_at", Table: TableMessages}

	ColMessageDeleted db.DBColumn = db.DBColumn{Column: "deleted", Table: TableMessages}
)
