package controller

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"monsoon/controller"
	"monsoon/middleware"
	"monsoon/mocks"

	"github.com/gin-gonic/gin"
	"go.uber.org/mock/gomock"
)

func TestMessageCreateRoute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		requestBody    string
		setupMocks     func(nats *mocks.MockINATSPublisher)
		expectedStatus int
	}{
		{
			name: "validation error: missing field",
			requestBody: `{
			}`,
			setupMocks: func(nats *mocks.MockINATSPublisher) {
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "successful message create",
			requestBody: `
			{
				"content": "hi guys"
			}
			`,
			setupMocks: func(nats *mocks.MockINATSPublisher) {
				nats.EXPECT().SendMsgNATS(gomock.Any()).Return()
			},
			expectedStatus: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockNATSPublisher := mocks.NewMockINATSPublisher(ctrl)
			if tt.setupMocks != nil {
				tt.setupMocks(mockNATSPublisher)
			}

			messageController := &controller.MessageController{NATS_PUB: mockNATSPublisher}

			gin.SetMode(gin.TestMode)
			r := gin.Default()

			r.POST("/message/create", middleware.RequireAuth(), messageController.MessageCreateRoute)

			req, _ := http.NewRequest("POST", "/message/create", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "12345")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}
