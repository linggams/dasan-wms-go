package service

import (
	"context"
	"fmt"

	"github.com/dppi/dppierp-api/internal/domain"
	"github.com/dppi/dppierp-api/internal/repository"
)

type CheckpointService struct {
	fabricRepo *repository.FabricRepository
	rackRepo   *repository.RackRepository
}

func NewCheckpointService(fabricRepo *repository.FabricRepository, rackRepo *repository.RackRepository) *CheckpointService {
	return &CheckpointService{
		fabricRepo: fabricRepo,
		rackRepo:   rackRepo,
	}
}

type ScanQRResponse struct {
	QRCode     string  `json:"qr_code"`
	Buyer      string  `json:"buyer"`
	Style      string  `json:"style"`
	Yard       string  `json:"yard"`
	QCResult   *string `json:"qc_result,omitempty"`
	FinishDate *string `json:"finish_date,omitempty"`
}

func (s *CheckpointService) GetOverview(ctx context.Context) ([]domain.MovementType, error) {
	return s.fabricRepo.GetMovementTypes(ctx)
}

func (s *CheckpointService) ScanQR(ctx context.Context, code string) (*ScanQRResponse, error) {
	fabric, err := s.fabricRepo.FindByCodeWithInventory(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("error finding fabric: %w", err)
	}
	if fabric == nil {
		return nil, fmt.Errorf("QR code is not found")
	}

	response := &ScanQRResponse{
		QRCode: fabric.Code,
		Buyer:  fabric.Buyer,
		Style:  fabric.Style,
		Yard:   fabric.Yard,
	}

	if fabric.Inventory != nil && fabric.Inventory.Stage == string(domain.StageQCFabric) {
		response.QCResult = fabric.QCResult
	}

	if fabric.Inventory != nil && fabric.Inventory.Stage == string(domain.StageRelaxation) {
		response.FinishDate = fabric.FinishDate
	}

	return response, nil
}

type MoveEntry struct {
	Code       string  `json:"code" binding:"required"`
	Yard       float64 `json:"yard,omitempty"`
	FinishDate string  `json:"finish_date,omitempty"`
	QCResult   string  `json:"qc_result,omitempty"`
}

type MoveRequest struct {
	Stage             string      `json:"stage" binding:"required"`
	BlockID           *int64      `json:"block_id,omitempty"`
	RackID            *int64      `json:"rack_id,omitempty"`
	RelaxationBlockID *int64      `json:"relaxation_block_id,omitempty"`
	RelaxationRackID  *int64      `json:"relaxation_rack_id,omitempty"`
	Entries           []MoveEntry `json:"entries" binding:"required,dive"`
}

func (s *CheckpointService) MoveStage(ctx context.Context, req *MoveRequest) error {
	if !domain.IsValidStage(req.Stage) {
		return fmt.Errorf("invalid stage: %s", req.Stage)
	}

	if len(req.Entries) == 0 {
		return fmt.Errorf("entries field is required")
	}

	repoReq := &repository.MoveRequestData{
		Stage:             req.Stage,
		BlockID:           req.BlockID,
		RackID:            req.RackID,
		RelaxationBlockID: req.RelaxationBlockID,
		RelaxationRackID:  req.RelaxationRackID,
		Entries:           make([]repository.MoveEntryData, len(req.Entries)),
	}

	for i, entry := range req.Entries {
		repoReq.Entries[i] = repository.MoveEntryData{
			Code:       entry.Code,
			Yard:       entry.Yard,
			FinishDate: entry.FinishDate,
			QCResult:   entry.QCResult,
		}
	}

	switch req.Stage {
	case string(domain.StageInventory):
		return s.fabricRepo.UpdateBlockRack(ctx, repoReq)
	case string(domain.StageRelaxation):
		return s.fabricRepo.UpdateRelaxationBlockRack(ctx, repoReq)
	case string(domain.StageQCFabric):
		return s.fabricRepo.UpdateStageWithQC(ctx, repoReq)
	default:
		return s.fabricRepo.UpdateStage(ctx, repoReq)
	}
}

type ScanRackResponse struct {
	Result  []ScanRackFabricItem `json:"result"`
	Summary RackSummary          `json:"summary"`
}

type ScanRackFabricItem struct {
	ID                int64   `json:"id"`
	Code              string  `json:"code"`
	Color             string  `json:"color,omitempty"`
	Lot               string  `json:"lot,omitempty"`
	Roll              string  `json:"roll,omitempty"`
	Weight            string  `json:"weight,omitempty"`
	Yard              string  `json:"yard"`
	FabricType        *string `json:"fabric_type,omitempty"`
	FabricContain     *string `json:"fabric_contain,omitempty"`
	FinishDate        *string `json:"finish_date,omitempty"`
	QCResult          *string `json:"qc_result,omitempty"`
	Buyer             string  `json:"buyer"`
	Style             string  `json:"style"`
	BlockID           *int64  `json:"block_id,omitempty"`
	RackID            *int64  `json:"rack_id,omitempty"`
	RelaxationBlockID *int64  `json:"relaxation_block_id,omitempty"`
	RelaxationRackID  *int64  `json:"relaxation_rack_id,omitempty"`
}

type RackSummary struct {
	TotalItems  int     `json:"total_items"`
	TotalYard   float64 `json:"total_yard"`
	TotalWeight float64 `json:"total_weight"`
	BlockName   string  `json:"block_name"`
	RackNumber  string  `json:"rack_number"`
}

func (s *CheckpointService) ScanRack(ctx context.Context, code string) (*ScanRackResponse, error) {
	rack, err := s.rackRepo.FindByName(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("error finding rack: %w", err)
	}
	if rack == nil {
		return nil, fmt.Errorf("rack not found")
	}

	fabrics, err := s.fabricRepo.GetFabricsByRackID(ctx, rack.ID)
	if err != nil {
		return nil, fmt.Errorf("error getting fabrics: %w", err)
	}

	var totalYard, totalWeight float64
	var blockName string
	var result []ScanRackFabricItem

	for _, f := range fabrics {
		var yard, weight float64
		fmt.Sscanf(f.Yard, "%f", &yard)
		fmt.Sscanf(f.Weight, "%f", &weight)
		totalYard += yard
		totalWeight += weight

		if f.Block != nil && blockName == "" {
			blockName = f.Block.Name
		}

		result = append(result, ScanRackFabricItem{
			ID:                f.ID,
			Code:              f.Code,
			Color:             f.Color,
			Lot:               f.Lot,
			Roll:              f.Roll,
			Weight:            f.Weight,
			Yard:              f.Yard,
			FabricType:        f.FabricType,
			FabricContain:     f.FabricContain,
			FinishDate:        f.FinishDate,
			QCResult:          f.QCResult,
			Buyer:             f.Buyer,
			Style:             f.Style,
			BlockID:           f.BlockID,
			RackID:            f.RackID,
			RelaxationBlockID: f.RelaxationBlockID,
			RelaxationRackID:  f.RelaxationRackID,
		})
	}

	if blockName == "" {
		blockName = "-"
	}

	return &ScanRackResponse{
		Result: result,
		Summary: RackSummary{
			TotalItems:  len(fabrics),
			TotalYard:   totalYard,
			TotalWeight: totalWeight,
			BlockName:   blockName,
			RackNumber:  rack.Name,
		},
	}, nil
}

type RelocationRequest struct {
	CurrentRackID int64 `json:"current_rack_id" binding:"required"`
	NewRackID     int64 `json:"new_rack_id" binding:"required"`
}

func (s *CheckpointService) Relocate(ctx context.Context, req *RelocationRequest) error {
	return s.fabricRepo.RelocateFabricsWithLog(ctx, req.CurrentRackID, req.NewRackID)
}
