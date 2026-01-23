package handler

import (
	"net/http"

	"github.com/dppi/dppierp-api/internal/service"
	"github.com/gin-gonic/gin"
)

type CheckpointHandler struct {
	service *service.CheckpointService
}

func NewCheckpointHandler(svc *service.CheckpointService) *CheckpointHandler {
	return &CheckpointHandler{service: svc}
}

func (h *CheckpointHandler) GetOverview(c *gin.Context) {
	stages, err := h.service.GetOverview(c.Request.Context())
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch overview.", err.Error())
		return
	}
	SuccessResponse(c, http.StatusOK, "Successfully fetched overview.", stages)
}

type ScanQRRequest struct {
	Code string `json:"code" binding:"required"`
}

func (h *CheckpointHandler) ScanQR(c *gin.Context) {
	var req ScanQRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ValidationErrorResponse(c, "The QR code is required.", map[string][]string{
			"code": {"The QR code is required."},
		})
		return
	}

	result, err := h.service.ScanQR(c.Request.Context(), req.Code)
	if err != nil {
		ErrorResponse(c, http.StatusNotFound, "Failed to founded QR.", err.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, "Successfully founded QR.", result)
}

type MoveStageRequest struct {
	BlockID           *int64              `json:"block_id,omitempty"`
	RackID            *int64              `json:"rack_id,omitempty"`
	RelaxationBlockID *int64              `json:"relaxation_block_id,omitempty"`
	RelaxationRackID  *int64              `json:"relaxation_rack_id,omitempty"`
	Entries           []service.MoveEntry `json:"entries" binding:"required,dive"`
}

func (h *CheckpointHandler) MoveStage(c *gin.Context) {
	stage := c.Query("stage")
	if stage == "" {
		ValidationErrorResponse(c, "Validation error.", map[string][]string{
			"stage": {"The stage field is required."},
		})
		return
	}

	var req MoveStageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ValidationErrorResponse(c, "Validation error.", map[string][]string{
			"entries": {"The entries field is required."},
		})
		return
	}

	if stage == "inventory" {
		if req.BlockID == nil {
			ValidationErrorResponse(c, "Block ID is required.", map[string][]string{
				"block_id": {"Block ID is required."},
			})
			return
		}
		if req.RackID == nil {
			ValidationErrorResponse(c, "Rack ID is required.", map[string][]string{
				"rack_id": {"Rack ID is required."},
			})
			return
		}
	}

	if stage == "relaxation" {
		if req.RelaxationBlockID == nil {
			ValidationErrorResponse(c, "Relaxation block ID is required.", map[string][]string{
				"relaxation_block_id": {"Relaxation block ID is required."},
			})
			return
		}
		if req.RelaxationRackID == nil {
			ValidationErrorResponse(c, "Relaxation rack ID is required.", map[string][]string{
				"relaxation_rack_id": {"Relaxation rack ID is required."},
			})
			return
		}
	}

	svcReq := &service.MoveRequest{
		Stage:             stage,
		BlockID:           req.BlockID,
		RackID:            req.RackID,
		RelaxationBlockID: req.RelaxationBlockID,
		RelaxationRackID:  req.RelaxationRackID,
		Entries:           req.Entries,
	}

	if err := h.service.MoveStage(c.Request.Context(), svcReq); err != nil {
		ErrorResponse(c, http.StatusUnprocessableEntity, "Failed to moved items.", err.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, "Successfully moved items.", true)
}

func (h *CheckpointHandler) ScanRack(c *gin.Context) {
	var req ScanQRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ValidationErrorResponse(c, "The rack code is required.", map[string][]string{
			"code": {"The rack code is required."},
		})
		return
	}

	result, err := h.service.ScanRack(c.Request.Context(), req.Code)
	if err != nil {
		ErrorResponse(c, http.StatusNotFound, "Failed to founded Rack QR.", err.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, "Successfully founded Rack QR.", result)
}

type RelocationRequest struct {
	CurrentRackID int64 `json:"current_rack_id" binding:"required"`
	NewRackID     int64 `json:"new_rack_id" binding:"required"`
}

func (h *CheckpointHandler) Relocate(c *gin.Context) {
	var req RelocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ValidationErrorResponse(c, "Invalid request body.", map[string][]string{
			"current_rack_id": {"The current rack id is required."},
			"new_rack_id":     {"The new rack id is required."},
		})
		return
	}

	if req.CurrentRackID == req.NewRackID {
		ValidationErrorResponse(c, "The current and new rack IDs must be different.", map[string][]string{
			"current_rack_id": {"The current and new rack IDs must be different."},
		})
		return
	}

	svcReq := &service.RelocationRequest{
		CurrentRackID: req.CurrentRackID,
		NewRackID:     req.NewRackID,
	}

	if err := h.service.Relocate(c.Request.Context(), svcReq); err != nil {
		ErrorResponse(c, http.StatusUnprocessableEntity, "Failed to relocated items.", err.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, "Successfully relocated items.", true)
}
