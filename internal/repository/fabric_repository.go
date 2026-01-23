package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/dppi/dppierp-api/internal/domain"
)

type FabricRepository struct {
	db *sql.DB
}

func NewFabricRepository(db *sql.DB) *FabricRepository {
	return &FabricRepository{db: db}
}

func (r *FabricRepository) GetMovementTypes(ctx context.Context) ([]domain.MovementType, error) {
	query := `SELECT id, name FROM movement_types ORDER BY id`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get movement types: %w", err)
	}
	defer rows.Close()

	var types []domain.MovementType
	for rows.Next() {
		var mt domain.MovementType
		if err := rows.Scan(&mt.ID, &mt.Name); err != nil {
			return nil, fmt.Errorf("failed to scan movement type: %w", err)
		}
		types = append(types, mt)
	}

	return types, nil
}

func (r *FabricRepository) FindByCode(ctx context.Context, code string) (*domain.Fabric, error) {
	query := `
		SELECT
			f.id, f.code, f.fabric_incoming_id, f.supplier_id, f.color, f.lot, f.roll,
			f.weight, f.width, f.yard, f.unit_id, f.fabric_type, f.fabric_contain,
			f.rack_id, f.block_id, f.relaxation_rack_id, f.relaxation_block_id,
			f.finish_date, f.qc_result, f.status, f.created_at, f.updated_at,
			COALESCE(b.name, '-') as buyer,
			COALESCE(o.style, '-') as style
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

func (r *FabricRepository) FindByCodeWithInventory(ctx context.Context, code string) (*domain.Fabric, error) {
	query := `
		SELECT
			f.id, f.code, f.fabric_incoming_id, f.supplier_id, f.color, f.lot, f.roll,
			f.weight, f.width, f.yard, f.unit_id, f.fabric_type, f.fabric_contain,
			f.rack_id, f.block_id, f.relaxation_rack_id, f.relaxation_block_id,
			f.finish_date, f.qc_result, f.status, f.created_at, f.updated_at,
			COALESCE(b.name, '-') as buyer,
			COALESCE(o.style, '-') as style,
			i.id as inv_id, i.stage as inv_stage
		FROM fabrics f
		LEFT JOIN fabric_incomings fi ON f.fabric_incoming_id = fi.id
		LEFT JOIN orders o ON fi.order_id = o.id
		LEFT JOIN buyers b ON o.buyer_id = b.id
		LEFT JOIN inventories i ON i.fabric_id = f.id AND i.deleted_at IS NULL
		WHERE f.code = ? AND f.deleted_at IS NULL
		LIMIT 1
	`

	var fabric domain.Fabric
	var finishDate, qcResult sql.NullString
	var invID sql.NullInt64
	var invStage sql.NullString

	err := r.db.QueryRowContext(ctx, query, code).Scan(
		&fabric.ID, &fabric.Code, &fabric.FabricIncomingID, &fabric.SupplierID,
		&fabric.Color, &fabric.Lot, &fabric.Roll, &fabric.Weight, &fabric.Width,
		&fabric.Yard, &fabric.UnitID, &fabric.FabricType, &fabric.FabricContain,
		&fabric.RackID, &fabric.BlockID, &fabric.RelaxationRackID, &fabric.RelaxationBlockID,
		&finishDate, &qcResult, &fabric.Status, &fabric.CreatedAt, &fabric.UpdatedAt,
		&fabric.Buyer, &fabric.Style,
		&invID, &invStage,
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

	if invID.Valid {
		fabric.Inventory = &domain.Inventory{
			ID:    invID.Int64,
			Stage: invStage.String,
		}
	}

	return &fabric, nil
}

func (r *FabricRepository) GetFabricsByRackID(ctx context.Context, rackID int64) ([]domain.Fabric, error) {
	query := `
		SELECT
			f.id, f.code, f.fabric_incoming_id, f.supplier_id, f.color, f.lot, f.roll,
			f.weight, f.width, f.yard, f.unit_id, f.fabric_type, f.fabric_contain,
			f.rack_id, f.block_id, f.relaxation_rack_id, f.relaxation_block_id,
			f.finish_date, f.qc_result, f.status, f.created_at, f.updated_at,
			COALESCE(b.name, '-') as buyer,
			COALESCE(o.style, '-') as style,
			blk.id as block_id_rel, blk.name as block_name
		FROM fabrics f
		LEFT JOIN fabric_incomings fi ON f.fabric_incoming_id = fi.id
		LEFT JOIN orders o ON fi.order_id = o.id
		LEFT JOIN buyers b ON o.buyer_id = b.id
		LEFT JOIN m_blocks blk ON f.block_id = blk.id
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
		var blockIDRel sql.NullInt64
		var blockName sql.NullString

		err := rows.Scan(
			&fabric.ID, &fabric.Code, &fabric.FabricIncomingID, &fabric.SupplierID,
			&fabric.Color, &fabric.Lot, &fabric.Roll, &fabric.Weight, &fabric.Width,
			&fabric.Yard, &fabric.UnitID, &fabric.FabricType, &fabric.FabricContain,
			&fabric.RackID, &fabric.BlockID, &fabric.RelaxationRackID, &fabric.RelaxationBlockID,
			&finishDate, &qcResult, &fabric.Status, &fabric.CreatedAt, &fabric.UpdatedAt,
			&fabric.Buyer, &fabric.Style,
			&blockIDRel, &blockName,
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

		if blockIDRel.Valid {
			fabric.Block = &domain.Block{
				ID:   blockIDRel.Int64,
				Name: blockName.String,
			}
		}

		fabrics = append(fabrics, fabric)
	}

	return fabrics, nil
}

type MoveRequestData struct {
	Stage             string
	BlockID           *int64
	RackID            *int64
	RelaxationBlockID *int64
	RelaxationRackID  *int64
	Entries           []MoveEntryData
}

type MoveEntryData struct {
	Code       string
	Yard       float64
	FinishDate string
	QCResult   string
}

func (r *FabricRepository) UpdateBlockRack(ctx context.Context, req *MoveRequestData) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	for _, entry := range req.Entries {
		code := entry.Code
		yard := entry.Yard

		fabric, err := r.FindByCodeWithInventory(ctx, code)
		if err != nil {
			return fmt.Errorf("error finding fabric %s: %w", code, err)
		}
		if fabric == nil {
			return fmt.Errorf("QR code %s is not found", code)
		}

		onStage := ""
		if fabric.Inventory != nil {
			onStage = fabric.Inventory.Stage
		}

		yardStr := strconv.FormatFloat(yard, 'f', -1, 64)
		if yard == 0 {
			yardStr = fabric.Yard
		}

		updateQuery := `UPDATE fabrics SET block_id = ?, rack_id = ?, yard = ?, updated_at = NOW() WHERE id = ? AND deleted_at IS NULL`
		_, err = tx.ExecContext(ctx, updateQuery, req.BlockID, req.RackID, yardStr, fabric.ID)
		if err != nil {
			return fmt.Errorf("failed to update fabric %s: %w", code, err)
		}

		storageQuery := `INSERT INTO fabric_storages (fabric_id, block_id, rack_id, created_by, created_at, updated_at) VALUES (?, ?, ?, ?, NOW(), NOW())`
		storageResult, err := tx.ExecContext(ctx, storageQuery, fabric.ID, req.BlockID, req.RackID, 1) // TODO: get auth user id
		if err != nil {
			return fmt.Errorf("failed to insert fabric storage log: %w", err)
		}
		storageLogID, _ := storageResult.LastInsertId()

		remarks := "From " + onStage
		if onStage == req.Stage {
			remarks = "Return " + onStage
		}

		invMovementID, err := r.handleStage(ctx, tx, fabric.ID, req.Stage, remarks, onStage)
		if err != nil {
			return fmt.Errorf("failed to handle stage: %w", err)
		}

		_, err = tx.ExecContext(ctx, `UPDATE fabric_storages SET inventory_movement_id = ? WHERE id = ?`, invMovementID, storageLogID)
		if err != nil {
			return fmt.Errorf("failed to update storage log: %w", err)
		}
	}

	return tx.Commit()
}

func (r *FabricRepository) UpdateRelaxationBlockRack(ctx context.Context, req *MoveRequestData) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	for _, entry := range req.Entries {
		code := entry.Code
		finishDate := entry.FinishDate

		fabric, err := r.FindByCodeWithInventory(ctx, code)
		if err != nil {
			return fmt.Errorf("error finding fabric %s: %w", code, err)
		}
		if fabric == nil {
			return fmt.Errorf("QR code %s is not found", code)
		}

		onStage := ""
		if fabric.Inventory != nil {
			onStage = fabric.Inventory.Stage
		}

		updateQuery := `UPDATE fabrics SET relaxation_block_id = ?, relaxation_rack_id = ?, finish_date = ?, updated_at = NOW() WHERE id = ? AND deleted_at IS NULL`
		_, err = tx.ExecContext(ctx, updateQuery, req.RelaxationBlockID, req.RelaxationRackID, finishDate, fabric.ID)
		if err != nil {
			return fmt.Errorf("failed to update fabric %s: %w", code, err)
		}

		relaxationQuery := `INSERT INTO fabric_relaxations (fabric_id, relaxation_block_id, relaxation_rack_id, finish_date, created_by, created_at, updated_at) VALUES (?, ?, ?, ?, ?, NOW(), NOW())`
		relaxationResult, err := tx.ExecContext(ctx, relaxationQuery, fabric.ID, req.RelaxationBlockID, req.RelaxationRackID, finishDate, 1) // TODO: get auth user id
		if err != nil {
			return fmt.Errorf("failed to insert fabric relaxation log: %w", err)
		}
		relaxationLogID, _ := relaxationResult.LastInsertId()

		remarks := "From " + onStage
		if onStage == req.Stage {
			remarks = "Return " + onStage
		}

		invMovementID, err := r.handleStage(ctx, tx, fabric.ID, req.Stage, remarks, onStage)
		if err != nil {
			return fmt.Errorf("failed to handle stage: %w", err)
		}

		_, err = tx.ExecContext(ctx, `UPDATE fabric_relaxations SET inventory_movement_id = ? WHERE id = ?`, invMovementID, relaxationLogID)
		if err != nil {
			return fmt.Errorf("failed to update relaxation log: %w", err)
		}
	}

	return tx.Commit()
}

func (r *FabricRepository) UpdateStageWithQC(ctx context.Context, req *MoveRequestData) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	for _, entry := range req.Entries {
		code := entry.Code
		qcResult := entry.QCResult
		if qcResult == "" {
			qcResult = "pass"
		}

		fabric, err := r.FindByCodeWithInventory(ctx, code)
		if err != nil {
			return fmt.Errorf("error finding fabric %s: %w", code, err)
		}
		if fabric == nil {
			return fmt.Errorf("QR code %s is not found", code)
		}

		onStage := ""
		if fabric.Inventory != nil {
			onStage = fabric.Inventory.Stage
		}

		updateQuery := `UPDATE fabrics SET qc_result = ?, updated_at = NOW() WHERE id = ? AND deleted_at IS NULL`
		_, err = tx.ExecContext(ctx, updateQuery, qcResult, fabric.ID)
		if err != nil {
			return fmt.Errorf("failed to update fabric %s: %w", code, err)
		}

		controlQuery := `INSERT INTO fabric_controls (fabric_id, result, created_by, created_at, updated_at) VALUES (?, ?, ?, NOW(), NOW())`
		controlResult, err := tx.ExecContext(ctx, controlQuery, fabric.ID, qcResult, 1) // TODO: get auth user id
		if err != nil {
			return fmt.Errorf("failed to insert fabric control log: %w", err)
		}
		controlLogID, _ := controlResult.LastInsertId()

		remarks := "From " + onStage
		if onStage == req.Stage {
			remarks = "Return " + onStage
		}

		invMovementID, err := r.handleStage(ctx, tx, fabric.ID, req.Stage, remarks, onStage)
		if err != nil {
			return fmt.Errorf("failed to handle stage: %w", err)
		}

		_, err = tx.ExecContext(ctx, `UPDATE fabric_controls SET inventory_movement_id = ? WHERE id = ?`, invMovementID, controlLogID)
		if err != nil {
			return fmt.Errorf("failed to update control log: %w", err)
		}
	}

	return tx.Commit()
}

func (r *FabricRepository) UpdateStage(ctx context.Context, req *MoveRequestData) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	for _, entry := range req.Entries {
		code := entry.Code

		fabric, err := r.FindByCodeWithInventory(ctx, code)
		if err != nil {
			return fmt.Errorf("error finding fabric %s: %w", code, err)
		}
		if fabric == nil {
			return fmt.Errorf("QR code %s is not found", code)
		}

		onStage := ""
		if fabric.Inventory != nil {
			onStage = fabric.Inventory.Stage
		}

		remarks := "From " + onStage
		if onStage == req.Stage {
			remarks = "Return " + onStage
		}

		_, err = r.handleStage(ctx, tx, fabric.ID, req.Stage, remarks, onStage)
		if err != nil {
			return fmt.Errorf("failed to handle stage: %w", err)
		}
	}

	return tx.Commit()
}

func (r *FabricRepository) handleStage(ctx context.Context, tx *sql.Tx, fabricID int64, toStage, remarks, onStage string) (int64, error) {
	now := time.Now()
	currentDate := now.Format("2006-01-02")
	currentTime := now.Format("15:04:05")

	// Get inventory_id for this fabric
	var inventoryID int64
	err := tx.QueryRowContext(ctx, `SELECT id FROM inventories WHERE fabric_id = ? AND deleted_at IS NULL`, fabricID).Scan(&inventoryID)
	if err != nil && err != sql.ErrNoRows {
		return 0, fmt.Errorf("failed to get inventory: %w", err)
	}

	// If no inventory exists, create one
	if err == sql.ErrNoRows {
		invResult, err := tx.ExecContext(ctx, `INSERT INTO inventories (fabric_id, stage, created_at, updated_at) VALUES (?, ?, ?, ?)`, fabricID, toStage, now, now)
		if err != nil {
			return 0, fmt.Errorf("failed to create inventory: %w", err)
		}
		inventoryID, _ = invResult.LastInsertId()
	}

	// Get movement_type_id based on stage
	var movementTypeID int64 = 1 // default
	_ = tx.QueryRowContext(ctx, `SELECT id FROM movement_types WHERE name = ? LIMIT 1`, toStage).Scan(&movementTypeID)

	// Get existing movement ID to update its time record
	var existingMovementID int64
	_ = tx.QueryRowContext(ctx, `SELECT id FROM inventory_movements WHERE fabric_id = ? AND status = 'starting' AND deleted_at IS NULL`, fabricID).Scan(&existingMovementID)

	// Finish existing movement
	_, err = tx.ExecContext(ctx, `UPDATE inventory_movements SET status = 'finished', updated_at = ? WHERE fabric_id = ? AND status = 'starting' AND deleted_at IS NULL`, now, fabricID)
	if err != nil {
		return 0, fmt.Errorf("failed to finish existing movement: %w", err)
	}

	// Update finish_date and finish_time in inventory_movement_times for the finished movement
	if existingMovementID > 0 {
		_, err = tx.ExecContext(ctx, `UPDATE inventory_movement_times SET finish_date = ?, finish_time = ?, updated_at = ? WHERE inventory_movement_id = ? AND deleted_at IS NULL`, currentDate, currentTime, now, existingMovementID)
		if err != nil {
			return 0, fmt.Errorf("failed to update movement time: %w", err)
		}
	}

	// Insert new movement
	movementQuery := `INSERT INTO inventory_movements (datetime, inventory_id, movement_type_id, fabric_id, remarks, status, action_by, created_at, updated_at) VALUES (?, ?, ?, ?, ?, 'starting', 1, ?, ?)`
	result, err := tx.ExecContext(ctx, movementQuery, now, inventoryID, movementTypeID, fabricID, remarks, now, now)
	if err != nil {
		return 0, fmt.Errorf("failed to insert inventory movement: %w", err)
	}
	invMovementID, _ := result.LastInsertId()

	// Insert movement time record
	timeQuery := `INSERT INTO inventory_movement_times (inventory_movement_id, start_date, start_time, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`
	_, err = tx.ExecContext(ctx, timeQuery, invMovementID, currentDate, currentTime, now, now)
	if err != nil {
		return 0, fmt.Errorf("failed to insert movement time: %w", err)
	}

	entryType := "out"
	if onStage == toStage {
		entryType = "actual"
	}

	entryQuery := `INSERT INTO inventory_entries (inventory_movement_id, type, from_stage, to_stage, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`
	_, err = tx.ExecContext(ctx, entryQuery, invMovementID, entryType, onStage, toStage, now, now)
	if err != nil {
		return 0, fmt.Errorf("failed to insert inventory entry: %w", err)
	}

	// Update inventory stage
	_, err = tx.ExecContext(ctx, `UPDATE inventories SET stage = ?, updated_at = ? WHERE id = ?`, toStage, now, inventoryID)
	if err != nil {
		return 0, fmt.Errorf("failed to update inventory stage: %w", err)
	}

	return invMovementID, nil
}

func (r *FabricRepository) RelocateFabricsWithLog(ctx context.Context, currentRackID, newRackID int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	fabricsQuery := `SELECT id FROM fabrics WHERE rack_id = ? AND deleted_at IS NULL`
	rows, err := tx.QueryContext(ctx, fabricsQuery, currentRackID)
	if err != nil {
		return fmt.Errorf("failed to get fabrics: %w", err)
	}

	var fabricIDs []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			rows.Close()
			return fmt.Errorf("failed to scan fabric id: %w", err)
		}
		fabricIDs = append(fabricIDs, id)
	}
	rows.Close()

	if len(fabricIDs) == 0 {
		return fmt.Errorf("no fabric found in the selected current rack")
	}

	now := time.Now()

	for _, fabricID := range fabricIDs {
		_, err = tx.ExecContext(ctx, `UPDATE fabrics SET rack_id = ?, updated_at = ? WHERE id = ?`, newRackID, now, fabricID)
		if err != nil {
			return fmt.Errorf("failed to update fabric rack: %w", err)
		}

		_, err = tx.ExecContext(ctx, `UPDATE fabric_rack_relocations SET is_archived = 1, updated_at = ? WHERE fabric_id = ? AND is_archived IS NULL`, now, fabricID)
		if err != nil {
			return fmt.Errorf("failed to archive relocation: %w", err)
		}

		relocationQuery := `INSERT INTO fabric_rack_relocations (fabric_id, current_rack_id, new_rack_id, created_by, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`
		_, err = tx.ExecContext(ctx, relocationQuery, fabricID, currentRackID, newRackID, 1, now, now) // TODO: get auth user id
		if err != nil {
			return fmt.Errorf("failed to insert relocation log: %w", err)
		}
	}

	return tx.Commit()
}

func (r *FabricRepository) UpdateFabricsForMove(ctx context.Context, codes []string, stage string, updates map[string]interface{}) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	for _, code := range codes {
		query := `UPDATE fabrics SET updated_at = NOW() WHERE code = ? AND deleted_at IS NULL`
		_, err := tx.ExecContext(ctx, query, code)
		if err != nil {
			return fmt.Errorf("failed to update fabric %s: %w", code, err)
		}

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
