package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dppi/dppierp-api/internal/domain"
	"github.com/dppi/dppierp-api/internal/middleware"
	"github.com/dppi/dppierp-api/internal/service"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// mockUserRepo is a simple mock for UserRepository
type mockUserRepo struct {
	user *domain.User
}

func (m *mockUserRepo) FindByEmail(email string) (*domain.User, error) {
	if m.user != nil && m.user.Email == email {
		return m.user, nil
	}
	return nil, nil
}

func TestAuthHandler_Login_Structure(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup Dependencies
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	mockRepo := &mockUserRepo{
		user: &domain.User{
			ID:              1,
			Name:            "Super Admin",
			Email:           "admin@dppi.com",
			Password:        string(hashedPassword),
			EmailVerifiedAt: &time.Time{},
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		},
	}
	authMiddleware := middleware.NewAuthMiddleware("secret-key-secret-key-secret-key-32")
	authService := service.NewAuthService(mockRepo, authMiddleware)
	authHandler := NewAuthHandler(authService)

	// Setup Router
	r := gin.New()
	r.POST("/auth/login", authHandler.Login)

	// Create Request
	reqBody := `{"email": "admin@dppi.com", "password": "password123"}`
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Execute
	r.ServeHTTP(w, req)

	// Assertions
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
		t.Log(w.Body.String())
	}

	var response struct {
		Status  string `json:"status"`
		Success bool   `json:"success"`
		Message string `json:"message"`
		Data    struct {
			TokenType       string `json:"token_type"`
			AccessToken     string `json:"access_token"`
			ExpiresIn       int    `json:"expires_in"`
			EmailVerifiedAt bool   `json:"email_verified_at"`
			UserInfo        struct {
				FullName string `json:"full_name"`
			} `json:"user_info"`
			HasAccess struct {
				Role string `json:"role"`
			} `json:"has_access"`
		} `json:"data"`
	}

	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}

	// Verify Structure Fields
	if response.Data.TokenType != "Bearer" {
		t.Errorf("Expected token_type 'Bearer', got '%s'", response.Data.TokenType)
	}
	if response.Data.AccessToken == "" {
		t.Error("Expected access_token to be present")
	}
	if response.Data.ExpiresIn != 3600 {
		t.Errorf("Expected expires_in 3600, got %d", response.Data.ExpiresIn)
	}
	if !response.Data.EmailVerifiedAt {
		t.Error("Expected email_verified_at to be true")
	}
	if response.Data.UserInfo.FullName != "Super Admin" {
		t.Errorf("Expected full_name 'Super Admin', got '%s'", response.Data.UserInfo.FullName)
	}
	if response.Data.HasAccess.Role != "superadmin" {
		t.Errorf("Expected role 'superadmin', got '%s'", response.Data.HasAccess.Role)
	}
}
