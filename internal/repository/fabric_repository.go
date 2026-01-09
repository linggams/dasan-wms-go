package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/dppi/dppierp-api/internal/domain"
)

type FabricRepository struct {
	db *sql.DB
}

func NewFabricRepository(db *sql.DB) *FabricRepository {
	return &FabricRepository{db: db}
}

// FindByCode finds fabric by QR code
func (r *FabricRepository) FindByCode(ctx context.Context, code string) (*domain.Fabric, error) {
	query := `
		SELECT
			f.id, f.code, f.fabric_incoming_id, f.supplier_id, f.color, f.lot, f.roll,
			f.weight, f.width, f.yard, f.unit_id, f.fabric_type, f.fabric_contain,
			f.rack_id, f.block_id, f.relaxation_rack_id, f.relaxation_block_id,
			f.finish_date, f.qc_result, f.status, f.created_at, f.updated_at,
			COALESCE(b.name, '-') as buyer,
			COALESCE(fi.style, '-') as style
		FROM fabrics f
		LEFT JOIN fabric_incomings fi ON f.fabric_incoming_id = fi.id
		LEFT JOIN orders o ON fi.order_id = o.id
		LEFT JOIN buyers b ON o.buyer_id = b.id
		WHERE f.code = ? AND f.deleted_at IS NULL
		LIMIT 1
	`

	var fabric domain.Fabric
	var finishDate sql.NullString
	var qcResult sql.NullString

	err := r.db.QueryRowContext(ctx, query, code).Scan(
		&fabric.ID, &fabric.Code, &fabric.FabricIncomingID, &fabric.SupplierID,
		&fabric.Color, &fabric.Lot, &fabric.Roll, &fabric.Weight, &fabric.Width,
		&fabric.Yard, &fabric.UnitID, &fabric.FabricType, &fabric.FabricContain,
		&fabric.RackID, &fabric.BlockID, &fabric.RelaxationRackID, &fabric.RelaxationBlockID,
		&finishDate, &qcResult, &fabric.Status, &fabric.CreatedAt, &fabric.UpdatedAt,
		&fabric.Buyer, &fabric.Style,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find fabric by code: %w", err)
	}

	if finishDate.Valid {
		fabric.FinishDate = &finishDate.String
	}
	if qcResult.Valid {
		fabric.QCResult = &qcResult.String
	}

	return &fabric, nil
}

// GetFabricsByRackID retrieves all fabrics in a rack
func (r *FabricRepository) GetFabricsByRackID(ctx context.Context, rackID int64) ([]domain.Fabric, error) {
	query := `
		SELECT
			f.id, f.code, f.fabric_incoming_id, f.supplier_id, f.color, f.lot, f.roll,
			f.weight, f.width, f.yard, f.unit_id, f.fabric_type, f.fabric_contain,
			f.rack_id, f.block_id, f.relaxation_rack_id, f.relaxation_block_id,
			f.finish_date, f.qc_result, f.status, f.created_at, f.updated_at,
			COALESCE(b.name, '-') as buyer,
			COALESCE(fi.style, '-') as style
		FROM fabrics f
		LEFT JOIN fabric_incomings fi ON f.fabric_incoming_id = fi.id
		LEFT JOIN orders o ON fi.order_id = o.id
		LEFT JOIN buyers b ON o.buyer_id = b.id
		WHERE f.rack_id = ? AND f.deleted_at IS NULL
	`

	rows, err := r.db.QueryContext(ctx, query, rackID)
	if err != nil {
		return nil, fmt.Errorf("failed to get fabrics by rack: %w", err)
	}
	defer rows.Close()

	var fabrics []domain.Fabric
	for rows.Next() {
		var fabric domain.Fabric
		var finishDate, qcResult sql.NullString

		err := rows.Scan(
			&fabric.ID, &fabric.Code, &fabric.FabricIncomingID, &fabric.SupplierID,
			&fabric.Color, &fabric.Lot, &fabric.Roll, &fabric.Weight, &fabric.Width,
			&fabric.Yard, &fabric.UnitID, &fabric.FabricType, &fabric.FabricContain,
			&fabric.RackID, &fabric.BlockID, &fabric.RelaxationRackID, &fabric.RelaxationBlockID,
			&finishDate, &qcResult, &fabric.Status, &fabric.CreatedAt, &fabric.UpdatedAt,
			&fabric.Buyer, &fabric.Style,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan fabric: %w", err)
		}

		if finishDate.Valid {
			fabric.FinishDate = &finishDate.String
		}
		if qcResult.Valid {
			fabric.QCResult = &qcResult.String
		}

		fabrics = append(fabrics, fabric)
	}

	return fabrics, nil
}

// UpdateFabricsForMove updates fabric records for stage movement
func (r *FabricRepository) UpdateFabricsForMove(ctx context.Context, codes []string, stage string, updates map[string]interface{}) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	for _, code := range codes {
		// Update fabric record based on stage
		query := `UPDATE fabrics SET updated_at = NOW() WHERE code = ? AND deleted_at IS NULL`
		_, err := tx.ExecContext(ctx, query, code)
		if err != nil {
			return fmt.Errorf("failed to update fabric %s: %w", code, err)
		}

		// Update or insert inventory record
		inventoryQuery := `
			UPDATE inventories
			SET stage = ?, updated_at = NOW()
			WHERE fabric_id = (SELECT id FROM fabrics WHERE code = ? AND deleted_at IS NULL)
			AND deleted_at IS NULL
		`
		_, err = tx.ExecContext(ctx, inventoryQuery, stage, code)
		if err != nil {
			return fmt.Errorf("failed to update inventory for fabric %s: %w", code, err)
		}
	}

	return tx.Commit()
}

// RelocateFabrics moves all fabrics from one rack to another
func (r *FabricRepository) RelocateFabrics(ctx context.Context, currentRackID, newRackID int64) error {
	query := `
		UPDATE fabrics
		SET rack_id = ?, updated_at = NOW()
		WHERE rack_id = ? AND deleted_at IS NULL
	`
	result, err := r.db.ExecContext(ctx, query, newRackID, currentRackID)
	if err != nil {
		return fmt.Errorf("failed to relocate fabrics: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("no fabric found in the selected current rack")
	}

	return nil
}
