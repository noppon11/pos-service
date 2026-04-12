package service

import (
	"context"
	"errors"
	"strings"
	"pos-service/internal/domain"
)

type DB interface {
	PingContext(ctx context.Context) error
}

type TenantValidator interface {
	TenantIDValidation(tenantID string) error
}

type BranchRepository interface {
	ListByTenantID(ctx context.Context, tenantID string) ([]domain.BranchResponse, error)
}

type PosService struct {
	db         DB
	branchRepo BranchRepository
	validator  TenantValidator
}

func NewPosService(db DB, branchRepo BranchRepository) *PosService {
	return &PosService{
		db:         db,
		branchRepo: branchRepo,
	}
}

func (s *PosService) GetHealth(ctx context.Context) error {
	return s.db.PingContext(ctx)
}

func (s *PosService) GetHealthByTenantID(ctx context.Context, tenantID string) error {
	if err := s.validator.TenantIDValidation(tenantID); err != nil {
		return err
	}
	return s.db.PingContext(ctx)
}

func (s *PosService) GetBranchesByTenantID(ctx context.Context, tenantID string) ([]domain.BranchResponse, error) {
	if err := s.validator.TenantIDValidation(tenantID); err != nil {
		return nil, err
	}

	branches, err := s.branchRepo.ListByTenantID(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	for _, branch := range branches {
		if err := validateBranch(branch); err != nil {
			return nil, err
		}
	}

	return branches, nil
}

func validateBranch(branch domain.BranchResponse) error {
	if strings.TrimSpace(branch.BranchID) == "" {
		return errors.New("branch_id is required")
	}
	if strings.TrimSpace(branch.BranchName) == "" {
		return errors.New("branch_name is required")
	}
	if branch.Status != "active" && branch.Status != "inactive" {
		return errors.New("status must be active or inactive")
	}
	return nil
}