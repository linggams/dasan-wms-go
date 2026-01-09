package service

import (
	"context"
	"testing"

	"github.com/dppi/dppierp-api/internal/domain"
)

// Mock FabricRepository for testing
type mockFabricRepository struct {
	fabrics map[string]*domain.Fabric
}

func newMockFabricRepository() *mockFabricRepository {
	return &mockFabricRepository{
		fabrics: map[string]*domain.Fabric{
			"F24120001": {
				ID:    1,
				Code:  "F24120001",
				Buyer: "GAP",
				Style: "Style-A",
				Yard:  "12.5",
			},
			"F24120002": {
				ID:    2,
				Code:  "F24120002",
				Buyer: "Target",
				Style: "Style-B",
				Yard:  "10.0",
			},
		},
	}
}

func (m *mockFabricRepository) FindByCode(ctx context.Context, code string) (*domain.Fabric, error) {
	if fabric, ok := m.fabrics[code]; ok {
		return fabric, nil
	}
	return nil, nil
}

func (m *mockFabricRepository) GetFabricsByRackID(ctx context.Context, rackID int64) ([]domain.Fabric, error) {
	return []domain.Fabric{}, nil
}

func (m *mockFabricRepository) UpdateFabricsForMove(ctx context.Context, codes []string, stage string, updates map[string]interface{}) error {
	return nil
}

func (m *mockFabricRepository) RelocateFabrics(ctx context.Context, currentRackID, newRackID int64) error {
	return nil
}

// Mock RackRepository for testing
type mockRackRepository struct {
	racks map[string]*domain.Rack
}

func newMockRackRepository() *mockRackRepository {
	return &mockRackRepository{
		racks: map[string]*domain.Rack{
			"RACK-001": {
				ID:   1,
				Name: "RACK-001",
			},
		},
	}
}

func (m *mockRackRepository) FindByName(ctx context.Context, name string) (*domain.Rack, error) {
	if rack, ok := m.racks[name]; ok {
		return rack, nil
	}
	return nil, nil
}

func (m *mockRackRepository) FindByID(ctx context.Context, id int64) (*domain.Rack, error) {
	for _, rack := range m.racks {
		if rack.ID == id {
			return rack, nil
		}
	}
	return nil, nil
}

func (m *mockRackRepository) GetBlockByID(ctx context.Context, id int64) (*domain.Block, error) {
	return nil, nil
}

func TestGetOverview(t *testing.T) {
	// Since GetOverview doesn't use repos, we can test it directly
	stages := domain.GetAllStages()

	if len(stages) != 9 {
		t.Errorf("Expected 9 stages, got %d", len(stages))
	}

	expectedStages := []string{
		"inventory", "cutting_wip", "stock_fabric", "cncm",
		"washing", "return_supplier", "destroy", "relaxation", "qc_fabric",
	}

	for i, stage := range stages {
		if stage.Name != expectedStages[i] {
			t.Errorf("Expected stage %s at position %d, got %s", expectedStages[i], i, stage.Name)
		}
	}
}

func TestIsValidStage(t *testing.T) {
	testCases := []struct {
		stage    string
		expected bool
	}{
		{"inventory", true},
		{"relaxation", true},
		{"cutting_wip", true},
		{"stock_fabric", true},
		{"cncm", true},
		{"washing", true},
		{"return_supplier", true},
		{"destroy", true},
		{"qc_fabric", true},
		{"invalid_stage", false},
		{"", false},
	}

	for _, tc := range testCases {
		result := domain.IsValidStage(tc.stage)
		if result != tc.expected {
			t.Errorf("IsValidStage(%s) = %v, expected %v", tc.stage, result, tc.expected)
		}
	}
}
