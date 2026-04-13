package service

import (
	"context"
	"errors"
	"strings"

	"pos-service/internal/domain"
)

var (
	ErrDBNotConfigured         = errors.New("database is not configured")
	ErrBranchRepoNotConfigured = errors.New("branch repository is not configured")
	ErrBranchIDRequired        = errors.New("branch_id is required")
	ErrBranchNameRequired      = errors.New("branch_name is required")
	ErrInvalidBranchStatus     = errors.New("status must be active or inactive")
)

type DB interface {
	PingContext(ctx context.Context) error
}

type NewInMemoryBranchRepository interface {
	ListByTenantID(ctx context.Context, tenantID string) ([]domain.BranchResponse, error)
	ByID(ctx context.Context, branchID string) (*domain.BranchResponse, error)
}

type PosService struct {
	db         DB
	branchRepo NewInMemoryBranchRepository
}

func NewPosService(db DB, branchRepo NewInMemoryBranchRepository) *PosService {
	return &PosService{
		db:         db,
		branchRepo: branchRepo,
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
		if err := validateBranch(branch); err != nil {
			return nil, err
		}
	}

	return branches, nil
}

func (s *PosService) GetBranchByID(ctx context.Context, branchID string) (*domain.BranchResponse, error) {
	if s.branchRepo == nil {
		return nil, ErrBranchRepoNotConfigured
	}

	branch, err := s.branchRepo.ByID(ctx, branchID)
	if err != nil {
		return nil, err
	}

	if branch == nil {
		return nil, nil
	}

	if err := validateBranch(*branch); err != nil {
		return nil, err
	}

	return branch, nil
}

func validateBranch(branch domain.BranchResponse) error {
	if strings.TrimSpace(branch.BranchID) == "" {
		return ErrBranchIDRequired
	}

	if strings.TrimSpace(branch.BranchName) == "" {
		return ErrBranchNameRequired
	}

	if !isValidBranchStatus(branch.Status) {
		return ErrInvalidBranchStatus
	}

	return nil
}

func isValidBranchStatus(status string) bool {
	switch status {
	case "active", "inactive":
		return true
	default:
		return false
	}
}