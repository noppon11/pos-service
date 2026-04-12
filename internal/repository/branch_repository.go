package repository

import (
	"context"

	"pos-service/internal/domain"
)

type BranchRepository interface {
	ListByTenantID(ctx context.Context, tenantID string) ([]domain.BranchResponse, error)
}