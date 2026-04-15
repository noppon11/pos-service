package repository

import (
	"context"
	"pos-service/internal/domain"
)

type ProductRepository interface {
	ListByTenantIDAndBranchID(ctx context.Context, tenantID string, branchID string) ([]domain.ProductResponse, error)
	GetByTenantIDBranchIDAndProductID(ctx context.Context, tenantID string, branchID string, productID string) (*domain.ProductResponse, error)
	Create(ctx context.Context, tenantID string, branchID string, product domain.ProductResponse) (*domain.ProductResponse, error)
	Update(ctx context.Context, tenantID string, branchID string, productID string, product domain.ProductResponse) (*domain.ProductResponse, error)
	Delete(ctx context.Context, tenantID string, branchID string, productID string) error
}