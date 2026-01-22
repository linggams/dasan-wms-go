package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/dppi/dppierp-api/internal/domain"
)

type RackRepository struct {
	db *sql.DB
}

func NewRackRepository(db *sql.DB) *RackRepository {
	return &RackRepository{db: db}
}

// FindByName finds rack by name/code
func (r *RackRepository) FindByName(ctx context.Context, name string) (*domain.Rack, error) {
	query := `
		SELECT id, name, created_at, updated_at
		FROM m_racks
		WHERE name = ? AND deleted_at IS NULL
		LIMIT 1
	`

	var rack domain.Rack
	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&rack.ID, &rack.Name, &rack.CreatedAt, &rack.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find rack by name: %w", err)
	}

	return &rack, nil
}

// FindByID finds rack by ID
func (r *RackRepository) FindByID(ctx context.Context, id int64) (*domain.Rack, error) {
	query := `
		SELECT id, name, created_at, updated_at
		FROM m_racks
		WHERE id = ? AND deleted_at IS NULL
		LIMIT 1
	`

	var rack domain.Rack
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&rack.ID, &rack.Name, &rack.CreatedAt, &rack.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find rack by id: %w", err)
	}

	return &rack, nil
}

// GetBlockByID finds block by ID
func (r *RackRepository) GetBlockByID(ctx context.Context, id int64) (*domain.Block, error) {
	query := `
		SELECT id, name, created_at, updated_at
		FROM m_blocks
		WHERE id = ? AND deleted_at IS NULL
		LIMIT 1
	`

	var block domain.Block
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&block.ID, &block.Name, &block.CreatedAt, &block.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find block by id: %w", err)
	}

	return &block, nil
}
