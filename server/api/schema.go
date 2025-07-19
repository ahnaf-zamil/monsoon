package api

/* Interfaces for API requests/responses */

type APIResponse struct {
	Err     bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	Status  int    `json:"status"`
}

type MessageCreateSchema struct {
	Content string `json:"content" binding:"required" example:"hunter2"`
}

type UserCreateSchema struct {
	Username    string `json:"username" binding:"required,min=3,max=15,alphanumunicode" example:"johndoe1"`
	DisplayName string `json:"display_name" binding:"required" example:"John Doe"`
	Email       string `json:"email" binding:"required,email" example:"john@doe.com"`

	Keys struct {
		Enc string `json:"enc" binding:"required"`
		Sig string `json:"sig" binding:"required"`
	} `json:"pub_keys" binding:"required"` // Base64 encoded X25519 and ED25519 public keys

	Salts struct {
		Enc  string `json:"enc_salt" binding:"required"`
		Auth string `json:"pw_salt" binding:"required"`
	}
	PasswordHash  string `json:"pw_hash" binding:"required"`
	EncryptedSeed string `json:"enc_seed" binding:"required"`
	Nonce         string `json:"nonce" binding:"required"`
}

type UserLoginSchema struct {
	Email    string `json:"email" binding:"required,email" example:"john@doe.com"`
	Password string `json:"password" binding:"required" example:"ilovejanedoe"`
}
