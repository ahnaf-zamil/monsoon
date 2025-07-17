package util

import (
	"github.com/gin-gonic/gin"
)

// Request validation helper

func ValidateRequestInput[T any](c *gin.Context, req *T) error {
	err := c.BindJSON(req)
	return err
}
