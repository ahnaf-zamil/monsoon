package lib

/* Data models for internal use and for JSON marshalling for API requests */

type MessageModel struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	CreatedAt int64  `json:"created_at"`
	UserID    string `json:"user_id"`
	RoomID    string `json:"room_id"`
}

type UserModel struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	CreatedAt   string `json:"created_at"`
	Email       string `json:"email"`
	Password    string `json:"-"`
}
