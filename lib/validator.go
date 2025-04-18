package lib

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Request validation helper

func ValidateRequestInput[T any](c *gin.Context, req *T) (*APIResponse, error) {
	rs := APIResponse{}
	if err := c.BindJSON(req); err != nil {

		rs.Err = true
		rs.Message = "Invalid input"
		c.JSON(http.StatusBadRequest, rs)
		return nil, err
	}

	return &rs, nil
}
