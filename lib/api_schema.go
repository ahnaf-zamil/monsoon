package lib

/* Interfaces for API requests/responses */

type APIResponse struct {
	Err     *bool           `json:"error"`
	Message *string         `json:"message"`
	Data    *map[string]any `json:"data"`
}
type MessageCreateSchema struct {
	Content string `json:"content" binding:"required"`
}
