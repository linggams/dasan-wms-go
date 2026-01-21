package handler

import (
	"net/http"

	"github.com/dppi/dppierp-api/internal/service"
	"github.com/gin-gonic/gin"
)

type MasterHandler struct {
	service *service.MasterService
}

func NewMasterHandler(service *service.MasterService) *MasterHandler {
	return &MasterHandler{service: service}
}

func (h *MasterHandler) GetBlocks(c *gin.Context) {
	blocks, err := h.service.GetAllBlocks()
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch blocks", err.Error())
		return
	}
	SuccessResponse(c, http.StatusOK, "Successfully fetched blocks", blocks)
}

func (h *MasterHandler) GetRacks(c *gin.Context) {
	racks, err := h.service.GetAllRacks()
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch racks", err.Error())
		return
	}
	SuccessResponse(c, http.StatusOK, "Successfully fetched racks", racks)
}

func (h *MasterHandler) GetRelaxationBlocks(c *gin.Context) {
	blocks, err := h.service.GetAllRelaxationBlocks()
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch relaxation blocks", err.Error())
		return
	}
	SuccessResponse(c, http.StatusOK, "Successfully fetched relaxation blocks", blocks)
}

func (h *MasterHandler) GetRelaxationRacks(c *gin.Context) {
	racks, err := h.service.GetAllRelaxationRacks()
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch relaxation racks", err.Error())
		return
	}
	SuccessResponse(c, http.StatusOK, "Successfully fetched relaxation racks", racks)
}
