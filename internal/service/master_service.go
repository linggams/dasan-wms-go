package service

import (
	"github.com/dppi/dppierp-api/internal/domain"
	"github.com/dppi/dppierp-api/internal/repository"
)

type MasterService struct {
	repo repository.MasterRepository
}

func NewMasterService(repo repository.MasterRepository) *MasterService {
	return &MasterService{repo: repo}
}

func (s *MasterService) GetAllBlocks() ([]domain.Block, error) {
	return s.repo.GetAllBlocks()
}

func (s *MasterService) GetAllRacks() ([]domain.Rack, error) {
	return s.repo.GetAllRacks()
}

func (s *MasterService) GetAllRelaxationBlocks() ([]domain.RelaxationBlock, error) {
	return s.repo.GetAllRelaxationBlocks()
}

func (s *MasterService) GetAllRelaxationRacks() ([]domain.RelaxationRack, error) {
	return s.repo.GetAllRelaxationRacks()
}
