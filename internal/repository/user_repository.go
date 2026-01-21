package repository

import (
	"database/sql"
	"errors"

	"github.com/dppi/dppierp-api/internal/domain"
)

type UserRepository interface {
	FindByEmail(email string) (*domain.User, error)
	FindByID(id int64) (*domain.User, error)
	UpdatePassword(userID int64, newHash string) error
	StoreResetToken(email, token string) error
	GetResetToken(email string) (string, error)
	DeleteResetToken(email string) error
}

type mysqlUserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &mysqlUserRepository{db: db}
}

func (r *mysqlUserRepository) FindByEmail(email string) (*domain.User, error) {
	query := `
		SELECT id, name, email, password, created_at, updated_at
		FROM users
		WHERE email = ? AND deleted_at IS NULL
	`

	user := &domain.User{}
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // User not found
		}
		return nil, err
	}

	return user, nil
}

func (r *mysqlUserRepository) FindByID(id int64) (*domain.User, error) {
	query := `
		SELECT id, name, email, password, created_at, updated_at
		FROM users
		WHERE id = ? AND deleted_at IS NULL
	`

	user := &domain.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // User not found
		}
		return nil, err
	}

	return user, nil
}

func (r *mysqlUserRepository) UpdatePassword(userID int64, newHash string) error {
	query := "UPDATE users SET password = ? WHERE id = ?"
	_, err := r.db.Exec(query, newHash, userID)
	return err
}

func (r *mysqlUserRepository) StoreResetToken(email, token string) error {
	// Delete existing token if any
	_, err := r.db.Exec("DELETE FROM password_reset_tokens WHERE email = ?", email)
	if err != nil {
		return err
	}

	query := "INSERT INTO password_reset_tokens (email, token, created_at) VALUES (?, ?, NOW())"
	_, err = r.db.Exec(query, email, token)
	return err
}

func (r *mysqlUserRepository) GetResetToken(email string) (string, error) {
	query := "SELECT token FROM password_reset_tokens WHERE email = ?"
	var token string
	err := r.db.QueryRow(query, email).Scan(&token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", err
	}
	return token, nil
}

func (r *mysqlUserRepository) DeleteResetToken(email string) error {
	_, err := r.db.Exec("DELETE FROM password_reset_tokens WHERE email = ?", email)
	return err
}
