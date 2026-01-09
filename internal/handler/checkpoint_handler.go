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

// GetOverview handles GET /check-point/v1/overview
func (h *CheckpointHandler) GetOverview(c *gin.Context) {
	stages := h.service.GetOverview(c.Request.Context())
	SuccessResponse(c, http.StatusOK, "Successfully fetched overview.", stages)
}

// ScanQRRequest represents the scan request body
type ScanQRRequest struct {
	Code string `json:"code" binding:"required"`
}

// ScanQR handles POST /check-point/v1/scan
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

// MoveStage handles POST /check-point/v1/move
func (h *CheckpointHandler) MoveStage(c *gin.Context) {
	stage := c.Query("stage")
	if stage == "" {
		ValidationErrorResponse(c, "The stage is required.", map[string][]string{
			"stage": {"The stage is required."},
		})
		return
	}

	var req service.MoveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ValidationErrorResponse(c, "The entries field is required.", map[string][]string{
			"entries": {"The entries field is required."},
		})
		return
	}

	req.Stage = stage

	if err := h.service.MoveStage(c.Request.Context(), &req); err != nil {
		ErrorResponse(c, http.StatusUnprocessableEntity, "Failed to moved items.", err.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, "Successfully moved items.", true)
}

// ScanRack handles POST /check-point/v1/scan-rack
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

// RelocationRequest represents the relocation request body
type RelocationRequest struct {
	CurrentRackID int64 `json:"current_rack_id" binding:"required"`
	NewRackID     int64 `json:"new_rack_id" binding:"required"`
}

// Relocate handles POST /check-point/v1/relocation
func (h *CheckpointHandler) Relocate(c *gin.Context) {
	var req RelocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ValidationErrorResponse(c, "Invalid request body.", map[string][]string{
			"current_rack_id": {"The current rack id is required."},
			"new_rack_id":     {"The new rack id is required."},
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
