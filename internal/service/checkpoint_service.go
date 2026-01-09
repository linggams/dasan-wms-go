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

// ScanQRResponse represents the scan response
type ScanQRResponse struct {
	QRCode     string  `json:"qr_code"`
	Buyer      string  `json:"buyer"`
	Style      string  `json:"style"`
	Yard       string  `json:"yard"`
	QCResult   *string `json:"qc_result,omitempty"`
	FinishDate *string `json:"finish_date,omitempty"`
}

// GetOverview returns all available stages
func (s *CheckpointService) GetOverview(ctx context.Context) []domain.StageInfo {
	return domain.GetAllStages()
}

// ScanQR scans a fabric QR code and returns its details
func (s *CheckpointService) ScanQR(ctx context.Context, code string) (*ScanQRResponse, error) {
	fabric, err := s.fabricRepo.FindByCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("error finding fabric: %w", err)
	}
	if fabric == nil {
		return nil, fmt.Errorf("QR code is not found")
	}

	response := &ScanQRResponse{
		QRCode:   fabric.Code,
		Buyer:    fabric.Buyer,
		Style:    fabric.Style,
		Yard:     fabric.Yard,
		QCResult: fabric.QCResult,
	}

	if fabric.FinishDate != nil {
		response.FinishDate = fabric.FinishDate
	}

	return response, nil
}

// MoveEntry represents an entry in the move request
type MoveEntry struct {
	Code       string  `json:"code" validate:"required"`
	Yard       float64 `json:"yard,omitempty"`
	FinishDate string  `json:"finish_date,omitempty"`
	QCResult   string  `json:"qc_result,omitempty"`
}

// MoveRequest represents the move stage request
type MoveRequest struct {
	Stage             string      `json:"stage" validate:"required"`
	BlockID           *int64      `json:"block_id,omitempty"`
	RackID            *int64      `json:"rack_id,omitempty"`
	RelaxationBlockID *int64      `json:"relaxation_block_id,omitempty"`
	RelaxationRackID  *int64      `json:"relaxation_rack_id,omitempty"`
	Entries           []MoveEntry `json:"entries" validate:"required,dive"`
}

// MoveStage moves fabrics to a new stage
func (s *CheckpointService) MoveStage(ctx context.Context, req *MoveRequest) error {
	if !domain.IsValidStage(req.Stage) {
		return fmt.Errorf("invalid stage: %s", req.Stage)
	}

	if len(req.Entries) == 0 {
		return fmt.Errorf("entries field is required")
	}

	codes := make([]string, len(req.Entries))
	for i, entry := range req.Entries {
		codes[i] = entry.Code
	}

	updates := make(map[string]interface{})
	if req.BlockID != nil {
		updates["block_id"] = *req.BlockID
	}
	if req.RackID != nil {
		updates["rack_id"] = *req.RackID
	}

	return s.fabricRepo.UpdateFabricsForMove(ctx, codes, req.Stage, updates)
}

// ScanRackResponse represents the scan rack response
type ScanRackResponse struct {
	Result  []domain.Fabric `json:"result"`
	Summary RackSummary     `json:"summary"`
}

type RackSummary struct {
	TotalItems  int     `json:"total_items"`
	TotalYard   float64 `json:"total_yard"`
	TotalWeight float64 `json:"total_weight"`
	BlockName   string  `json:"block_name"`
	RackNumber  string  `json:"rack_number"`
}

// ScanRack scans a rack and returns all fabrics in it
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

	// Calculate summary
	var totalYard, totalWeight float64
	var blockName string
	for _, f := range fabrics {
		var yard, weight float64
		fmt.Sscanf(f.Yard, "%f", &yard)
		fmt.Sscanf(f.Weight, "%f", &weight)
		totalYard += yard
		totalWeight += weight

		if f.Block != nil && blockName == "" {
			blockName = f.Block.Name
		}
	}

	return &ScanRackResponse{
		Result: fabrics,
		Summary: RackSummary{
			TotalItems:  len(fabrics),
			TotalYard:   totalYard,
			TotalWeight: totalWeight,
			BlockName:   blockName,
			RackNumber:  rack.Name,
		},
	}, nil
}

// RelocationRequest represents the relocation request
type RelocationRequest struct {
	CurrentRackID int64 `json:"current_rack_id" validate:"required"`
	NewRackID     int64 `json:"new_rack_id" validate:"required"`
}

// Relocate moves all fabrics from one rack to another
func (s *CheckpointService) Relocate(ctx context.Context, req *RelocationRequest) error {
	return s.fabricRepo.RelocateFabrics(ctx, req.CurrentRackID, req.NewRackID)
}
