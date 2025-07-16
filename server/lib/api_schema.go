package lib

/* Interfaces for API requests/responses */

type APIResponse struct {
	Err     bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type MessageCreateSchema struct {
	Content string `json:"content" binding:"required" example:"hunter2"`
}

type UserCreateSchema struct {
	Username    string `json:"username" binding:"required,min=3,max=15,alphanumunicode" example:"johndoe1"`
	DisplayName string `json:"display_name" binding:"required" example:"John Doe"`
	Email       string `json:"email" binding:"required,email" example:"john@doe.com"`
	Password    string `json:"password" binding:"required,min=8" example:"ilovejanedoe"`
}

type UserLoginSchema struct {
	Email    string `json:"email" binding:"required,email" example:"john@doe.com"`
	Password string `json:"password" binding:"required" example:"ilovejanedoe"`
}
