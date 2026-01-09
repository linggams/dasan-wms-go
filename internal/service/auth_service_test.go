package service

import (
	"testing"
	"time"

	"github.com/dppi/dppierp-api/internal/domain"
	"github.com/dppi/dppierp-api/internal/middleware"
	"golang.org/x/crypto/bcrypt"
)

// MockUserRepository
type mockUserRepository struct {
	users map[string]*domain.User
}

func (m *mockUserRepository) FindByEmail(email string) (*domain.User, error) {
	if user, ok := m.users[email]; ok {
		return user, nil
	}
	return nil, nil // Not found
}

func TestAuthService_Login(t *testing.T) {
	// Setup
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	mockRepo := &mockUserRepository{
		users: map[string]*domain.User{
			"test@example.com": {
				ID:        1,
				Name:      "Test User",
				Email:     "test@example.com",
				Password:  string(hashedPassword),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
	}

	authMiddleware := middleware.NewAuthMiddleware("secret-key-secret-key-secret-key-32")
	authService := NewAuthService(mockRepo, authMiddleware)

	// Test Case 1: Success
	token, refreshToken, user, err := authService.Login("test@example.com", password)
	if err != nil {
		t.Errorf("Expected success, got error: %v", err)
	}
	if token == "" {
		t.Error("Expected access token to be generated")
	}
	if refreshToken == "" {
		t.Error("Expected refresh token to be generated")
	}
	if user.Email != "test@example.com" {
		t.Errorf("Expected user email test@example.com, got %v", user.Email)
	}

	// Test Case 2: Wrong Password
	_, _, _, err = authService.Login("test@example.com", "wrongpassword")
	if err == nil {
		t.Error("Expected error for wrong password, got nil")
	}
	if err.Error() != "invalid credentials" {
		t.Errorf("Expected 'invalid credentials', got '%v'", err.Error())
	}

	// Test Case 3: User Not Found
	_, _, _, err = authService.Login("unknown@example.com", password)
	if err == nil {
		t.Error("Expected error for unknown user, got nil")
	}
	if err.Error() != "invalid credentials" {
		t.Errorf("Expected 'invalid credentials', got '%v'", err.Error())
	}
}
