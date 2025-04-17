package lib

/* Interfaces for API requests/responses */

type APIResponse struct {
	Err     bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}
type MessageCreateSchema struct {
	Content string `json:"content" binding:"required"`
}

type UserCreateSchema struct {
	Username    string `json:"username" binding:"required,min=3,max=15,alphanumunicode"`
	DisplayName string `json:"display_name" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=8"`
}
