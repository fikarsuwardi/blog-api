package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"blog-api/internal/database"
	"blog-api/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) {
	var err error
	database.DB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Migrate tables
	database.DB.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{})
}

func TestRegister(t *testing.T) {
	setupTestDB(t)

	tests := []struct {
		name           string
		requestBody    RegisterRequest
		expectedStatus int
		expectError    bool
	}{
		{
			name: "Valid Registration",
			requestBody: RegisterRequest{
				Email:    "test@example.com",
				Password: "password123",
				Name:     "Test User",
			},
			expectedStatus: http.StatusCreated,
			expectError:    false,
		},
		{
			name: "Missing Email",
			requestBody: RegisterRequest{
				Email:    "",
				Password: "password123",
				Name:     "Test User",
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "Short Password",
			requestBody: RegisterRequest{
				Email:    "test2@example.com",
				Password: "123",
				Name:     "Test User",
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "Invalid Email Format",
			requestBody: RegisterRequest{
				Email:    "invalid-email",
				Password: "password123",
				Name:     "Test User",
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/api/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			Register(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			var response map[string]interface{}
			json.NewDecoder(w.Body).Decode(&response)

			if tt.expectError {
				if _, exists := response["error"]; !exists {
					t.Error("Expected error in response")
				}
			} else {
				if _, exists := response["token"]; !exists {
					t.Error("Expected token in response")
				}
			}
		})
	}
}

func TestLogin(t *testing.T) {
	setupTestDB(t)

	// Create test user first
	registerBody := RegisterRequest{
		Email:    "login@example.com",
		Password: "password123",
		Name:     "Login Test",
	}
	body, _ := json.Marshal(registerBody)
	req := httptest.NewRequest("POST", "/api/register", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	Register(w, req)

	tests := []struct {
		name           string
		requestBody    LoginRequest
		expectedStatus int
		expectError    bool
	}{
		{
			name: "Valid Login",
			requestBody: LoginRequest{
				Email:    "login@example.com",
				Password: "password123",
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name: "Invalid Password",
			requestBody: LoginRequest{
				Email:    "login@example.com",
				Password: "wrongpassword",
			},
			expectedStatus: http.StatusUnauthorized,
			expectError:    true,
		},
		{
			name: "Non-existent User",
			requestBody: LoginRequest{
				Email:    "notfound@example.com",
				Password: "password123",
			},
			expectedStatus: http.StatusUnauthorized,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/api/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			Login(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			var response map[string]interface{}
			json.NewDecoder(w.Body).Decode(&response)

			if tt.expectError {
				if _, exists := response["error"]; !exists {
					t.Error("Expected error in response")
				}
			} else {
				if _, exists := response["token"]; !exists {
					t.Error("Expected token in response")
				}
			}
		})
	}
}
