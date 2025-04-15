package lib

import "time"

/* Data models for internal use */

type Message struct {
	ID        string
	Content   string
	CreatedAt time.Time
	UserID    string
	RoomID    string
}
