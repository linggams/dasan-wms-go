package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"

	"github.com/dppi/dppierp-api/internal/domain"
	"github.com/dppi/dppierp-api/internal/middleware"
	"github.com/dppi/dppierp-api/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo       repository.UserRepository
	authMiddleware *middleware.AuthMiddleware
}

func NewAuthService(userRepo repository.UserRepository, authMiddleware *middleware.AuthMiddleware) *AuthService {
	return &AuthService{
		userRepo:       userRepo,
		authMiddleware: authMiddleware,
	}
}

// Login authenticates a user
func (s *AuthService) Login(username, password string) (string, string, *domain.User, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return "", "", nil, err
	}
	if user == nil {
		return "", "", nil, errors.New("invalid credentials")
	}

	// Verify password (bcrypt)
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", "", nil, errors.New("invalid credentials")
	}

	// Generate Token
	accessToken, err := s.authMiddleware.GenerateToken(user.ID, user.Email, user.Name)
	if err != nil {
		return "", "", nil, err
	}

	refreshToken, err := s.authMiddleware.GenerateRefreshToken(user.ID)
	if err != nil {
		return "", "", nil, err
	}

	return accessToken, refreshToken, user, nil
}

// RefreshToken validates the refresh token and returns a new access/refresh token pair
func (s *AuthService) RefreshToken(tokenString string) (string, string, error) {
	userID, err := s.authMiddleware.ValidateRefreshToken(tokenString)
	if err != nil {
		return "", "", err
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return "", "", err
	}
	if user == nil {
		return "", "", errors.New("user not found")
	}

	newAccessToken, err := s.authMiddleware.GenerateToken(user.ID, user.Email, user.Name)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := s.authMiddleware.GenerateRefreshToken(user.ID)
	if err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}

// Logout performs any necessary server-side logout operations
func (s *AuthService) Logout() error {
	return nil
}

// ForgotPasswordRequest handles the request to send a password reset link
func (s *AuthService) ForgotPasswordRequest(email string) error {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("email not found")
	}

	// Generate a random token
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return err
	}
	token := hex.EncodeToString(b)

	// Store token in DB
	if err := s.userRepo.StoreResetToken(email, token); err != nil {
		return err
	}

	return nil
}

// ResetPassword resets the user's password using the token
func (s *AuthService) ResetPassword(email, token, newPassword string) error {
	storedToken, err := s.userRepo.GetResetToken(email)
	if err != nil {
		return err
	}
	if storedToken == "" || storedToken != token {
		return errors.New("invalid token")
	}

	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	hashedPassword, err := HashPassword(newPassword)
	if err != nil {
		return err
	}

	if err := s.userRepo.UpdatePassword(user.ID, hashedPassword); err != nil {
		return err
	}

	// Delete used token
	return s.userRepo.DeleteResetToken(email)
}

// ChangePassword allows an authenticated user to change their password
func (s *AuthService) ChangePassword(userID int64, oldPassword, newPassword string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return errors.New("invalid old password")
	}

	hashedPassword, err := HashPassword(newPassword)
	if err != nil {
		return err
	}

	return s.userRepo.UpdatePassword(user.ID, hashedPassword)
}

// HashPassword is a utility to hash passwords
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
