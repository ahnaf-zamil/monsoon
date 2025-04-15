package lib

/* Data models for internal use */

type MessageModel struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	CreatedAt int64  `json:"created_at"`
	UserID    string `json:"user_id"`
	RoomID    string `json:"room_id"`
}
