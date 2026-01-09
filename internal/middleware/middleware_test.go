package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestAuthMiddleware_GenerateAndValidateToken(t *testing.T) {
	auth := NewAuthMiddleware("test-secret-key-12345")

	// Generate token
	token, err := auth.GenerateToken(1, "test@example.com", "Test User")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	if token == "" {
		t.Error("Expected non-empty token")
	}

	// Validate token
	claims, err := auth.ValidateToken(token)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	if claims.UserID != 1 {
		t.Errorf("Expected UserID 1, got %d", claims.UserID)
	}

	if claims.Email != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got '%s'", claims.Email)
	}

	if claims.Name != "Test User" {
		t.Errorf("Expected name 'Test User', got '%s'", claims.Name)
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	auth := NewAuthMiddleware("test-secret-key-12345")

	// Try to validate invalid token
	_, err := auth.ValidateToken("invalid-token")
	if err == nil {
		t.Error("Expected error for invalid token")
	}
}

func TestAuthMiddleware_ExpiredToken(t *testing.T) {
	// This test would require a way to create expired tokens
	// Skipping for now as it requires time manipulation
	t.Skip("Requires time manipulation to test expired tokens")
}

func TestAuthMiddleware_MissingAuthHeader(t *testing.T) {
	auth := NewAuthMiddleware("test-secret-key-12345")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("GET", "/test", nil)
	// No Authorization header

	middleware := auth.Authenticate()
	middleware(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestAuthMiddleware_InvalidAuthFormat(t *testing.T) {
	auth := NewAuthMiddleware("test-secret-key-12345")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("GET", "/test", nil)
	c.Request.Header.Set("Authorization", "InvalidFormat")

	middleware := auth.Authenticate()
	middleware(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestAuthMiddleware_ValidAuth(t *testing.T) {
	auth := NewAuthMiddleware("test-secret-key-12345")
	w := httptest.NewRecorder()
	c, router := gin.CreateTestContext(w)

	// Generate a valid token
	token, _ := auth.GenerateToken(1, "test@example.com", "Test User")

	c.Request, _ = http.NewRequest("GET", "/test", nil)
	c.Request.Header.Set("Authorization", "Bearer "+token)

	router.Use(auth.Authenticate())
	router.GET("/test", func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		c.JSON(http.StatusOK, gin.H{"user_id": userID})
	})

	router.ServeHTTP(w, c.Request)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestCORSMiddleware(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("OPTIONS", "/test", nil)
	c.Request.Header.Set("Origin", "http://example.com")

	middleware := CORSMiddleware("*")
	middleware(c)

	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Error("Expected CORS header to be set")
	}

	if w.Code != 204 {
		t.Errorf("Expected status 204 for OPTIONS, got %d", w.Code)
	}
}

func TestLogger(t *testing.T) {
	w := httptest.NewRecorder()
	c, router := gin.CreateTestContext(w)

	router.Use(Logger())
	router.GET("/test", func(c *gin.Context) {
		time.Sleep(10 * time.Millisecond) // Small delay for latency
		c.Status(http.StatusOK)
	})

	c.Request, _ = http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, c.Request)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}
