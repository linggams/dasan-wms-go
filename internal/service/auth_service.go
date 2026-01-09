package service

import (
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

// Login authenticates a user by email and password
func (s *AuthService) Login(email, password string) (string, string, *domain.User, error) {
	// 1. Fetch user by email
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", "", nil, err
	}
	if user == nil {
		return "", "", nil, errors.New("invalid credentials")
	}

	// 2. Verify password (bcrypt)
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", "", nil, errors.New("invalid credentials")
	}

	// 3. Generate Token
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

// HashPassword is a utility to hash passwords (useful for registration or seeding)
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
