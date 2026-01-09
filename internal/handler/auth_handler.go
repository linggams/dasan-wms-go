package handler

import (
	"net/http"

	"github.com/dppi/dppierp-api/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// LoginRequest represents the login request body
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	TokenType       string `json:"token_type"`
	AccessToken     string `json:"access_token"`
	RefreshToken    string `json:"refresh_token"`
	ExpiresIn       int    `json:"expires_in"`
	EmailVerifiedAt bool   `json:"email_verified_at"`
	UserInfo        struct {
		PhotoPath string `json:"photo_path"`
		FullName  string `json:"full_name"`
	} `json:"user_info"`
	HasAccess struct {
		UserID int64  `json:"user_id"`
		Role   string `json:"role"`
	} `json:"has_access"`
}

// Login handles POST /auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ValidationErrorResponse(c, "Validation Error", map[string][]string{
			"email":    {"The email field is required and must be valid."},
			"password": {"The password field is required with minimum 6 characters."},
		})
		return
	}

	token, refreshToken, user, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Username or Password is wrong!",
			"errors":  nil,
		})
		return
	}

	response := LoginResponse{
		TokenType:       "Bearer",
		AccessToken:     token,
		RefreshToken:    refreshToken,
		ExpiresIn:       3600, // Default 1 hour
		EmailVerifiedAt: user.EmailVerifiedAt != nil,
	}

	// User Info
	response.UserInfo.FullName = user.Name
	if user.ProfilePhotoPath != nil {
		response.UserInfo.PhotoPath = *user.ProfilePhotoPath
	} else {
		response.UserInfo.PhotoPath = "http://127.0.0.1:8000/images/avatar.png"
	}

	response.HasAccess.UserID = user.ID
	response.HasAccess.Role = "superadmin"

	SuccessResponse(c, http.StatusOK, "Successfully logged.", response)
}

// Me handles GET /auth/me
func (h *AuthHandler) Me(c *gin.Context) {
	userID, _ := c.Get("user_id")
	email, _ := c.Get("email")
	name, _ := c.Get("name")

	SuccessResponse(c, http.StatusOK, "User fetched successfully.", gin.H{
		"id":    userID,
		"email": email,
		"name":  name,
	})
}
