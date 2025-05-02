package controller

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"ws_realtime_app/controller"
	"ws_realtime_app/lib"
	"ws_realtime_app/mocks"

	"github.com/gin-gonic/gin"
	"go.uber.org/mock/gomock"
)

func TestUserCreateRoute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		requestBody    string
		setupMocks     func(db *mocks.MockIUserDB, hasher *mocks.MockIPasswordHasher)
		expectedStatus int
	}{
		{
			name: "validation error: missing field",
			requestBody: `{
				"display_name": "User 1",
				"email": "user1@example.com",
			}`,
			setupMocks: func(db *mocks.MockIUserDB, hasher *mocks.MockIPasswordHasher) {
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "validation error: invalid username",
			requestBody: `{
				"username": "user @ 1",
				"display_name": "User 1",
				"email": "user1@example.com",
				"password": "password123"
			}`,
			setupMocks: func(db *mocks.MockIUserDB, hasher *mocks.MockIPasswordHasher) {
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "validation error: invalid email",
			requestBody: `{
				"username": "user1",
				"display_name": "User 1",
				"email": "user1example.com",
				"password": "password123"
			}`,
			setupMocks: func(db *mocks.MockIUserDB, hasher *mocks.MockIPasswordHasher) {
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "validation error: invalid password",
			requestBody: `{
				"username": "user1",
				"display_name": "User 1",
				"email": "user1@example.com",
				"password": "pas"
			}`, // Too short password
			setupMocks: func(db *mocks.MockIUserDB, hasher *mocks.MockIPasswordHasher) {
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "successful user creation",
			requestBody: `{
				"username": "user1",
				"display_name": "User 1",
				"email": "user1@example.com",
				"password": "password123"
			}`,
			setupMocks: func(db *mocks.MockIUserDB, hasher *mocks.MockIPasswordHasher) {
				db.EXPECT().GetUserByAnyField(gomock.Any(), gomock.Any()).Return(nil, nil)
				hasher.EXPECT().Hash(gomock.Any()).Return([]byte{}, nil)
				db.EXPECT().CreateUser(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "user already exists",
			requestBody: `{
				"username": "user1",
				"display_name": "User1",
				"email": "user1@example.com",
				"password": "password123"
			}`,
			setupMocks: func(db *mocks.MockIUserDB, hasher *mocks.MockIPasswordHasher) {
				db.EXPECT().GetUserByAnyField(gomock.Any(), gomock.Any()).
					Return(&lib.UserModel{Username: "user1"}, nil)
			},
			expectedStatus: http.StatusConflict,
		},
		{
			name: "database error",
			requestBody: `{
				"username": "dberror",
				"display_name": "Oops",
				"email": "fail@example.com",
				"password": "password123"
			}`,
			setupMocks: func(db *mocks.MockIUserDB, hasher *mocks.MockIPasswordHasher) {
				db.EXPECT().GetUserByAnyField(gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf("db failed"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "hasher error",
			requestBody: `{
				"username": "user1",
				"display_name": "User 1",
				"email": "user1@example.com",
				"password": "password123"
			}`,
			setupMocks: func(db *mocks.MockIUserDB, hasher *mocks.MockIPasswordHasher) {
				db.EXPECT().GetUserByAnyField(gomock.Any(), gomock.Any()).Return(nil, nil)
				hasher.EXPECT().Hash(gomock.Any()).Return(nil, fmt.Errorf("hasher error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserDB := mocks.NewMockIUserDB(ctrl)
			mockPasswordHasher := mocks.NewMockIPasswordHasher(ctrl)
			if tt.setupMocks != nil {
				tt.setupMocks(mockUserDB, mockPasswordHasher)
			}

			userController := &controller.UserController{UserDB: mockUserDB, PasswordHasher: mockPasswordHasher}
			gin.SetMode(gin.TestMode)
			r := gin.Default()
			r.POST("/user/create", userController.UserCreateRoute)

			req, _ := http.NewRequest("POST", "/user/create", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestUserLoginRoute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		requestBody    string
		setupMocks     func(db *mocks.MockIUserDB, hasher *mocks.MockIPasswordHasher)
		expectedStatus int
	}{
		{
			name: "validation error: invalid email",
			requestBody: `{
				"email": "user1example.com",
				"password": "password123"
			}`,
			setupMocks: func(db *mocks.MockIUserDB, hasher *mocks.MockIPasswordHasher) {
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "successful user login",
			requestBody: `
			{
			  "email": "user1@example.com",
			  "password": "password123"
			}
			`,
			setupMocks: func(db *mocks.MockIUserDB, hasher *mocks.MockIPasswordHasher) {
				db.EXPECT().GetUserByAnyField(gomock.Any(), gomock.Any()).Return(&lib.UserModel{Username: "user1"}, nil)
				hasher.EXPECT().Verify(gomock.Any(), gomock.Any()).Return(true, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "unsuccessful user login",
			requestBody: `
			{
			  "email": "user1@example.com",
			  "password": "wrongpassword"
			}
			`,
			setupMocks: func(db *mocks.MockIUserDB, hasher *mocks.MockIPasswordHasher) {
				db.EXPECT().GetUserByAnyField(gomock.Any(), gomock.Any()).Return(&lib.UserModel{Username: "user1"}, nil)
				hasher.EXPECT().Verify(gomock.Any(), gomock.Any()).Return(false, nil)
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "user doesn't exist",
			requestBody: `{
				"email": "user1@example.com",
				"password": "password123"
			}`,
			setupMocks: func(db *mocks.MockIUserDB, hasher *mocks.MockIPasswordHasher) {
				db.EXPECT().GetUserByAnyField(gomock.Any(), gomock.Any()).
					Return(nil, nil)
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "database error",
			requestBody: `{
				"email": "fail@example.com",
				"password": "password123"
			}`,
			setupMocks: func(db *mocks.MockIUserDB, hasher *mocks.MockIPasswordHasher) {
				db.EXPECT().GetUserByAnyField(gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf("db failed"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "hasher error",
			requestBody: `{
				"email": "user1@example.com",
				"password": "password123"
			}`,
			setupMocks: func(db *mocks.MockIUserDB, hasher *mocks.MockIPasswordHasher) {
				db.EXPECT().GetUserByAnyField(gomock.Any(), gomock.Any()).Return(&lib.UserModel{Username: "user1"}, nil)
				hasher.EXPECT().Verify(gomock.Any(), gomock.Any()).Return(false, fmt.Errorf("hasher error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserDB := mocks.NewMockIUserDB(ctrl)
			mockPasswordHasher := mocks.NewMockIPasswordHasher(ctrl)

			if tt.setupMocks != nil {
				tt.setupMocks(mockUserDB, mockPasswordHasher)
			}

			userController := &controller.UserController{UserDB: mockUserDB, PasswordHasher: mockPasswordHasher}

			gin.SetMode(gin.TestMode)
			r := gin.Default()
			r.POST("/user/login", userController.UserLoginRoute)

			req, _ := http.NewRequest("POST", "/user/login", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("%s", w.Body)
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}
