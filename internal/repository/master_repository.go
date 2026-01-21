package repository

import (
	"database/sql"

	"github.com/dppi/dppierp-api/internal/domain"
)

type MasterRepository interface {
	GetAllBlocks() ([]domain.Block, error)
	GetAllRacks() ([]domain.Rack, error)
	GetAllRelaxationBlocks() ([]domain.RelaxationBlock, error)
	GetAllRelaxationRacks() ([]domain.RelaxationRack, error)
}

type mysqlMasterRepository struct {
	db *sql.DB
}

func NewMasterRepository(db *sql.DB) MasterRepository {
	return &mysqlMasterRepository{db: db}
}

func (r *mysqlMasterRepository) GetAllBlocks() ([]domain.Block, error) {
	query := "SELECT id, name, created_at, updated_at, deleted_at FROM blocks WHERE deleted_at IS NULL"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blocks []domain.Block
	for rows.Next() {
		var b domain.Block
		if err := rows.Scan(&b.ID, &b.Name, &b.CreatedAt, &b.UpdatedAt, &b.DeletedAt); err != nil {
			return nil, err
		}
		blocks = append(blocks, b)
	}
	return blocks, nil
}

func (r *mysqlMasterRepository) GetAllRacks() ([]domain.Rack, error) {
	query := "SELECT id, name, created_at, updated_at, deleted_at FROM racks WHERE deleted_at IS NULL"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var racks []domain.Rack
	for rows.Next() {
		var ra domain.Rack
		if err := rows.Scan(&ra.ID, &ra.Name, &ra.CreatedAt, &ra.UpdatedAt, &ra.DeletedAt); err != nil {
			return nil, err
		}
		racks = append(racks, ra)
	}
	return racks, nil
}

func (r *mysqlMasterRepository) GetAllRelaxationBlocks() ([]domain.RelaxationBlock, error) {
	query := "SELECT id, name, created_at, updated_at, deleted_at FROM m_relaxation_blocks WHERE deleted_at IS NULL"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blocks []domain.RelaxationBlock
	for rows.Next() {
		var b domain.RelaxationBlock
		if err := rows.Scan(&b.ID, &b.Name, &b.CreatedAt, &b.UpdatedAt, &b.DeletedAt); err != nil {
			return nil, err
		}
		blocks = append(blocks, b)
	}
	return blocks, nil
}

func (r *mysqlMasterRepository) GetAllRelaxationRacks() ([]domain.RelaxationRack, error) {
	query := "SELECT id, name, created_at, updated_at, deleted_at FROM m_relaxation_racks WHERE deleted_at IS NULL"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var racks []domain.RelaxationRack
	for rows.Next() {
		var ra domain.RelaxationRack
		if err := rows.Scan(&ra.ID, &ra.Name, &ra.CreatedAt, &ra.UpdatedAt, &ra.DeletedAt); err != nil {
			return nil, err
		}
		racks = append(racks, ra)
	}
	return racks, nil
}
