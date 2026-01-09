package domain

import "testing"

func TestGetAllStages(t *testing.T) {
	stages := GetAllStages()

	if len(stages) != 9 {
		t.Errorf("Expected 9 stages, got %d", len(stages))
	}

	// Verify stage IDs are sequential
	for i, stage := range stages {
		if stage.ID != i+1 {
			t.Errorf("Expected stage ID %d, got %d", i+1, stage.ID)
		}
	}
}

func TestIsValidStage(t *testing.T) {
	validStages := []string{
		"inventory",
		"relaxation",
		"cutting_wip",
		"stock_fabric",
		"cncm",
		"washing",
		"return_supplier",
		"destroy",
		"qc_fabric",
	}

	for _, stage := range validStages {
		if !IsValidStage(stage) {
			t.Errorf("Expected stage '%s' to be valid", stage)
		}
	}

	invalidStages := []string{
		"invalid",
		"",
		"INVENTORY",
		"unknown_stage",
	}

	for _, stage := range invalidStages {
		if IsValidStage(stage) {
			t.Errorf("Expected stage '%s' to be invalid", stage)
		}
	}
}

func TestStageConstants(t *testing.T) {
	if StageInventory != "inventory" {
		t.Errorf("Expected StageInventory to be 'inventory', got '%s'", StageInventory)
	}

	if StageRelaxation != "relaxation" {
		t.Errorf("Expected StageRelaxation to be 'relaxation', got '%s'", StageRelaxation)
	}

	if StageCuttingWIP != "cutting_wip" {
		t.Errorf("Expected StageCuttingWIP to be 'cutting_wip', got '%s'", StageCuttingWIP)
	}

	if StageQCFabric != "qc_fabric" {
		t.Errorf("Expected StageQCFabric to be 'qc_fabric', got '%s'", StageQCFabric)
	}
}
