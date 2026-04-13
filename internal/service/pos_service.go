package service

import (
	"context"
	"errors"

	"pos-service/internal/domain"
)

var (
	ErrDBNotConfigured         = errors.New("database is not configured")
	ErrBranchRepoNotConfigured = errors.New("branch repository is not configured")
	ErrValidatorNotConfigured  = errors.New("validator not configured")
	ErrInvalidBranchStatus           = errors.New("status must be active or inactive")
	ErrInvalidBranchCurrency         = errors.New("currency must be 3 uppercase letters")
	ErrInvalidBranchCurrencyRequired = errors.New("currency is required")
	ErrInvalidBranchTimezoneRequired = errors.New("timezone is required")
	ErrBranchIDRequired              = errors.New("branch_id is required")
	ErrBranchNameRequired            = errors.New("branch_name is required")
)

type DB interface {
	PingContext(ctx context.Context) error
}

type NewInMemoryBranchRepository interface {
	ListByTenantID(ctx context.Context, tenantID string) ([]domain.BranchResponse, error)
	GetByTenantIDAndBranchID(ctx context.Context, tenantID string, branchID string) (*domain.BranchResponse, error)
}

type PosService struct {
	db         DB
	branchRepo NewInMemoryBranchRepository
	validator  Validator
}

type Validator interface {
	BranchValidation(branch domain.BranchResponse) error
}

func NewPosService(db DB, branchRepo NewInMemoryBranchRepository, v Validator) *PosService {
	return &PosService{
		db:         db,
		branchRepo: branchRepo,
		validator:  v,
	}
}

func (s *PosService) GetHealth(ctx context.Context) error {
	if s.db == nil {
		return ErrDBNotConfigured
	}

	return s.db.PingContext(ctx)
}

func (s *PosService) GetHealthByTenantID(ctx context.Context, tenantID string) error {
	if s.db == nil {
		return ErrDBNotConfigured
	}

	return s.db.PingContext(ctx)
}

func (s *PosService) GetBranchesByTenantID(ctx context.Context, tenantID string) ([]domain.BranchResponse, error) {
	if s.branchRepo == nil {
		return nil, ErrBranchRepoNotConfigured
	}

	branches, err := s.branchRepo.ListByTenantID(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	for _, branch := range branches {
		if err := s.validateBranch(branch); err != nil {
			return nil, err
		}
	}

	return branches, nil
}

func (s *PosService) GetBranchDetail(ctx context.Context, tenantID, branchID string) (*domain.BranchResponse, error) {
	if s.branchRepo == nil {
		return nil, ErrBranchRepoNotConfigured
	}

	branch, err := s.branchRepo.GetByTenantIDAndBranchID(ctx, tenantID, branchID)
	if err != nil {
		return nil, err
	}

	if branch == nil {
		return nil, nil
	}

	if err := s.validateBranch(*branch); err != nil {
		return nil, err
	}

	return branch, nil
}

func (s *PosService) validateBranch(branch domain.BranchResponse) error {

	if s.validator == nil {
		return ErrValidatorNotConfigured
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
