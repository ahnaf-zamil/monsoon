package api

/* Data models for internal use and for JSON marshalling for API requests */

type MessageModel struct {
	ID             string `db:"id" json:"id"`
	Content        string `db:"content" json:"content"`
	EditedAt       any    `db:"edited_at" json:"edited_at"`
	CreatedAt      int64  `db:"created_at" json:"created_at"`
	AuthorID       string `db:"author_id" json:"author_id"`
	ConversationID string `db:"conversation_id" json:"conversation_id,omitempty"` // Optional
	Deleted        bool   `json:"-"`
}

type UserModel struct {
	ID           string `db:"id" json:"id"`
	Username     string `db:"username" json:"username"`
	DisplayName  string `db:"display_name" json:"display_name"`
	CreatedAt    string `db:"created_at" json:"created_at"`
	Email        string `db:"-" json:"email,omitempty"`
	PasswordHash []byte `db:"-" json:"-"`
}

type UserSessionModel struct {
	SessionID    string
	UserID       string
	RefreshToken string
	CreatedAt    int64
}

type UserAuthModel struct {
	UserID         string `json:"-"`
	Email          string `json:"-"`
	PasswordHash   []byte `json:"-"`
	PasswordSalt   []byte `json:"-"`
	EncryptionSalt []byte `json:"enc_salt"`
	EncryptedSeed  []byte `json:"enc_seed"`
	Nonce          []byte `json:"nonce"`
}

type DirectConversationModel struct {
	ConversationID string
	User1ID        string
	User2ID        string
	CreatedAt      int64
}

type ConversationParticipant struct {
	ConversationID string `db:"conversation_id" json:"conversation_id"`
	UserID         string `db:"user_id" json:"-"`
	JoinedAt       int64  `db:"joined_at"`
	Role           string `db:"role"`
}

type InboxConversation struct {
	ConversationID string `db:"conversation_id" json:"conversation_id"`
	Type           string `db:"type" json:"type"`
	Name           any    `db:"name" json:"name"`
	UpdatedAt      int64  `db:"updated_at" json:"updated_at"`
	UserID         any    `db:"user_id" json:"user_id,omitempty"`
}
