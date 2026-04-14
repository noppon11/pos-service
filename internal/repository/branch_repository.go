package repository

import (
	"context"
	"pos-service/internal/domain"
)

type BranchRepository interface {
	ListByTenantID(ctx context.Context, tenantID string) ([]domain.BranchResponse, error)
	GetByTenantIDAndBranchID(ctx context.Context, tenantID, branchID string) (*domain.BranchResponse, error)
}