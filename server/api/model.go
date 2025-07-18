package api

/* Data models for internal use and for JSON marshalling for API requests */

type MessageModel struct {
	ID          string `json:"id"`
	Content     string `json:"content"`
	CreatedAt   int64  `json:"created_at"`
	AuthorID    string `json:"author_id"`
	RecipientID string `json:"recipient_id,omitempty"` // Optional
	RoomID      string `json:"room_id,omitempty"`      // Optional
	IsDM        bool   `json:"is_dm"`
}

type UserModel struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	CreatedAt   string `json:"created_at"`
	Email       string `json:"email"`
	Password    string `json:"-"`
}

type UserSessionModel struct {
	SessionID    string
	UserID       string
	RefreshToken string
	CreatedAt    int64
}
