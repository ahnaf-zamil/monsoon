package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"monsoon/middleware"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

func TestRequireAuth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name            string
		requestBody     string
		setupMocks      func()
		expectedStatus  int
		extraReqHeaders [][2]string
	}{
		{
			name: "auth success",
			requestBody: `{
			}`,
			setupMocks: func() {
			},
			expectedStatus: http.StatusOK,
			extraReqHeaders: [][2]string{
				{"Authorization", "12345"},
			},
		},
		{
			name: "auth failure",
			requestBody: `{
			}`,
			setupMocks: func() {
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMocks != nil {
				tt.setupMocks()
			}

			gin.SetMode(gin.TestMode)
			r := gin.Default()

			r.POST("/auth_required", middleware.RequireAuth(), func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{})
			})

			req, _ := http.NewRequest("POST", "/auth_required", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")

			for _, x := range tt.extraReqHeaders {
				req.Header.Set(x[0], x[1])
			}

			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}
