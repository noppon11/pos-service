package service

import (
	"context"
	"database/sql"

	"pos-service/internal/domain"
	appErr "pos-service/internal/errors"
)

type BranchRepository interface {
	ListByTenantID(ctx context.Context, tenantID string) ([]domain.BranchResponse, error)
	GetByTenantIDAndBranchID(ctx context.Context, tenantID string, branchID string) (*domain.BranchResponse, error)
}

type Validator interface {
	BranchValidation(branch domain.BranchResponse) error
}

type PosService struct {
	db          *sql.DB
	branchRepo BranchRepository
	validator  Validator
}

func NewPosService(db  *sql.DB, branchRepo BranchRepository, v Validator) *PosService {
	return &PosService{
		db:         db,
		branchRepo: branchRepo,
		validator:  v,
	}
}

func (s *PosService) GetHealth(ctx context.Context) error {
	if s.db == nil {
		return appErr.ErrDBNotConfigured
	}

	return s.db.PingContext(ctx)
}

func (s *PosService) GetHealthByTenantID(ctx context.Context, tenantID string) error {
	if s.db == nil {
		return appErr.ErrDBNotConfigured
	}

	return s.db.PingContext(ctx)
}

func (s *PosService) GetBranchesByTenantID(ctx context.Context, tenantID string) ([]domain.BranchResponse, error) {
	if s.branchRepo == nil {
		return nil, appErr.ErrBranchRepoNotConfigured
	}

	branches, err := s.branchRepo.ListByTenantID(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	for i := range branches {
		if err := s.validateBranch(branches[i]); err != nil {
			return nil, err
		}
	}

	return branches, nil
}

func (s *PosService) GetBranchDetail(ctx context.Context, tenantID, branchID string) (*domain.BranchResponse, error) {
	if s.branchRepo == nil {
		return nil, appErr.ErrBranchRepoNotConfigured
	}

	branch, err := s.branchRepo.GetByTenantIDAndBranchID(ctx, tenantID, branchID)
	if err != nil {
		return nil, err
	}

	if err := s.validateBranch(*branch); err != nil {
		return nil, err
	}

	return branch, nil
}

func (s *PosService) validateBranch(branch domain.BranchResponse) error {
	if s.validator == nil {
		return appErr.ErrValidatorNotConfigured
	}

	return s.validator.BranchValidation(branch)
}

func isValidBranchStatus(status string) bool {
	switch status {
	case "active", "inactive":
		return true
	default:
		return false
	}
}
