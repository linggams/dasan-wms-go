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

// RefreshTokenRequest represents the refresh token request body
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// ForgotPasswordRequest represents the forgot password request body
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResetPasswordRequest represents the reset password request body
type ResetPasswordRequest struct {
	Email                string `json:"email" binding:"required,email"`
	Token                string `json:"token" binding:"required"`
	Password             string `json:"password" binding:"required,min=6"`
	PasswordConfirmation string `json:"password_confirmation" binding:"required,eqfield=Password"`
}

// ChangePasswordRequest represents the change password request body
type ChangePasswordRequest struct {
	CurrentPassword      string `json:"current_password" binding:"required"`
	Password             string `json:"password" binding:"required,min=6"`
	PasswordConfirmation string `json:"password_confirmation" binding:"required,eqfield=Password"`
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

// RefreshToken handles POST /auth/token/refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ValidationErrorResponse(c, "Refresh token is required.", map[string][]string{
			"refresh_token": {"The refresh token field is required."},
		})
		return
	}

	newToken, newRefreshToken, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Invalid refresh token",
		})
		return
	}

	SuccessResponse(c, http.StatusOK, "Token refreshed successfully.", gin.H{
		"access_token":  newToken,
		"refresh_token": newRefreshToken,
		"token_type":    "Bearer",
		"expires_in":    3600,
	})
}

// Logout handles POST /auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	h.authService.Logout()
	SuccessResponse(c, http.StatusOK, "Successfully logged out.", nil)
}

// ForgotPasswordRequest handles POST /auth/forgot-password/request
func (h *AuthHandler) ForgotPasswordRequest(c *gin.Context) {
	var req ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ValidationErrorResponse(c, "Email is required.", map[string][]string{
			"email": {"The email field is required."},
		})
		return
	}

	if err := h.authService.ForgotPasswordRequest(req.Email); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Failed to process request", err.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, "Reset link sent to your email.", nil)
}

// ResetPassword handles POST /auth/forgot-password/reset
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ValidationErrorResponse(c, "Validation Error", map[string][]string{
			"email":    {"The email field is required."},
			"token":    {"The token field is required."},
			"password": {"The password field is required."},
		})
		return
	}

	if err := h.authService.ResetPassword(req.Email, req.Token, req.Password); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Failed to reset password", err.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, "Password reset successfully.", nil)
}

// ChangePassword handles POST /profile/change-password
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ValidationErrorResponse(c, "Validation Error", map[string][]string{
			"current_password": {"The current password field is required."},
			"password":         {"The password field is required."},
		})
		return
	}

	if err := h.authService.ChangePassword(userID.(int64), req.CurrentPassword, req.Password); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Failed to change password", err.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, "Password changed successfully.", nil)
}
