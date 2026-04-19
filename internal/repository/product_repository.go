package repository

import (
	"context"
	"pos-service/internal/domain"
)

type ProductListFilter struct {
	Page       int
	Limit      int
	CategoryID string
}

type ProductRepository interface {
	ListByTenantIDAndBranchID(
		ctx context.Context,
		tenantID string,
		branchID string,
		filter ProductListFilter,
	) ([]domain.Product, int, error)

	GetByTenantIDBranchIDAndProductID(
		ctx context.Context,
		tenantID string,
		branchID string,
		productID string,
	) (*domain.Product, error)

	Create(
		ctx context.Context,
		tenantID string,
		branchID string,
		product domain.Product,
	) (*domain.Product, error)

	Update(
		ctx context.Context,
		tenantID string,
		branchID string,
		productID string,
		product domain.Product,
	) (*domain.Product, error)

	Delete(
		ctx context.Context,
		tenantID string,
		branchID string,
		productID string,
	) error
}