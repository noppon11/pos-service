package service

import (
	"context"
	"database/sql"

	"pos-service/internal/domain"
	"pos-service/internal/dto"
	appErr "pos-service/internal/errors"
)

type BranchRepository interface {
	ListByTenantID(ctx context.Context, tenantID string) ([]domain.BranchResponse, error)
	GetByTenantIDAndBranchID(ctx context.Context, tenantID string, branchID string) (*domain.BranchResponse, error)
}

type ProductRepository interface {
	ListByTenantIDAndBranchID(ctx context.Context, tenantID string, branchID string) ([]domain.ProductResponse, error)
	GetByTenantIDBranchIDAndProductID(ctx context.Context, tenantID string, branchID string, productID string) (*domain.ProductResponse, error)
	Create(ctx context.Context, tenantID string, branchID string, product domain.ProductResponse) (*domain.ProductResponse, error)
	Update(ctx context.Context, tenantID string, branchID string, productID string, product domain.ProductResponse) (*domain.ProductResponse, error)
	Delete(ctx context.Context, tenantID string, branchID string, productID string) error
}

type Validator interface {
	ValidateBranch(branch domain.BranchResponse) error
	ValidateProduct(product domain.ProductResponse) error
}

type PosService struct {
	db          *sql.DB
	branchRepo  BranchRepository
	productRepo ProductRepository
	validator   Validator
}

func NewPosService(
	db *sql.DB,
	branchRepo BranchRepository,
	productRepo ProductRepository,
	v Validator,
) *PosService {
	return &PosService{
		db:          db,
		branchRepo:  branchRepo,
		productRepo: productRepo,
		validator:   v,
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

	// For now this is still DB-level health only.
	// If you want tenant-specific health, add a repository existence check here.
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

	if branch == nil {
		return nil, appErr.ErrBranchNotFound
	}

	if err := s.validateBranch(*branch); err != nil {
		return nil, err
	}

	return branch, nil
}

func (s *PosService) GetProducts(ctx context.Context, tenantID, branchID string) ([]domain.ProductResponse, error) {
	if s.productRepo == nil {
		return nil, appErr.ErrProductRepoNotConfigured
	}

	products, err := s.productRepo.ListByTenantIDAndBranchID(ctx, tenantID, branchID)
	if err != nil {
		return nil, err
	}

	for i := range products {
		if err := s.validateProduct(products[i]); err != nil {
			return nil, err
		}
	}

	return products, nil
}

func (s *PosService) GetProductByID(ctx context.Context, tenantID, branchID, productID string) (*domain.ProductResponse, error) {
	if s.productRepo == nil {
		return nil, appErr.ErrProductRepoNotConfigured
	}

	product, err := s.productRepo.GetByTenantIDBranchIDAndProductID(ctx, tenantID, branchID, productID)
	if err != nil {
		return nil, err
	}

	if product == nil {
		return nil, appErr.ErrProductNotFound
	}

	if err := s.validateProduct(*product); err != nil {
		return nil, err
	}

	return product, nil
}

func (s *PosService) CreateNewProduct(
	ctx context.Context,
	tenantID, branchID string,
	req dto.CreateProductRequest,
) (*domain.ProductResponse, error) {
	if err := s.ensureProductDependencies(); err != nil {
		return nil, err
	}

	if tenantID == "" {
		return nil, appErr.ErrTenantIDRequired
	}
	if branchID == "" {
		return nil, appErr.ErrBranchIDRequired
	}

	product := domain.ProductResponse{
		Name:       req.Name,
		SKU:        req.SKU,
		Price:      req.Price,
		CategoryID: req.CategoryID,
		Unit:       req.Unit,
		IsActive:   req.IsActive,
	}

	if err := s.validateProduct(product); err != nil {
		return nil, err
	}

	created, err := s.productRepo.Create(ctx, tenantID, branchID, product)
	if err != nil {
		return nil, err
	}

	if created == nil {
		return nil, appErr.ErrCreateProductFailed
	}

	if err := s.validateProduct(*created); err != nil {
		return nil, err
	}

	return created, nil
}
func (s *PosService) UpdateProduct(
	ctx context.Context,
	tenantID, branchID, productID string,
	req dto.UpdateProductRequest,
) (*domain.ProductResponse, error) {
	if err := s.ensureProductDependencies(); err != nil {
		return nil, err
	}

	if productID == "" {
		return nil, appErr.ErrProductIDRequired
	}

	if err := s.validateProductScope(ctx, tenantID, branchID); err != nil {
		return nil, err
	}

	product := domain.ProductResponse{
		Name:       req.Name,
		SKU:        req.SKU,
		Price:      req.Price,
		CategoryID: req.CategoryID,
		Unit:       req.Unit,
		IsActive:   req.IsActive,
	}

	if err := s.validateProduct(product); err != nil {
		return nil, err
	}

	updated, err := s.productRepo.Update(ctx, tenantID, branchID, productID, product)
	if err != nil {
		return nil, err
	}

	if updated == nil {
		return nil, appErr.ErrProductNotFound
	}

	if err := s.validateProduct(*updated); err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *PosService) DeleteProduct(ctx context.Context, tenantID, branchID, productID string) error {
	if err := s.ensureProductDependencies(); err != nil {
		return err
	}

	if err := s.validateProductScope(ctx, tenantID, branchID); err != nil {
		return err
	}

	product, err := s.productRepo.GetByTenantIDBranchIDAndProductID(ctx, tenantID, branchID, productID)
	if err != nil {
		return err
	}
	if product == nil {
		return appErr.ErrProductNotFound
	}

	return s.productRepo.Delete(ctx, tenantID, branchID, productID)
}

func (s *PosService) validateBranch(branch domain.BranchResponse) error {
	if s.validator == nil {
		return appErr.ErrValidatorNotConfigured
	}

	return s.validator.ValidateBranch(branch)
}

func (s *PosService) validateProduct(product domain.ProductResponse) error {
	if s.validator == nil {
		return appErr.ErrValidatorNotConfigured
	}

	return s.validator.ValidateProduct(product)
}

func (s *PosService) ensureProductDependencies() error {
	if s.productRepo == nil {
		return appErr.ErrProductRepoNotConfigured
	}

	if s.branchRepo == nil {
		return appErr.ErrBranchRepoNotConfigured
	}

	return nil
}

func (s *PosService) validateProductScope(
	ctx context.Context,
	tenantID, branchID string,
) error {
	branch, err := s.branchRepo.GetByTenantIDAndBranchID(ctx, tenantID, branchID)
	if err != nil {
		return err
	}
	if branch == nil {
		return appErr.ErrBranchNotFound
	}

	return nil
}

func validateRequiredIDs(tenantID, branchID, productID string) error {
	if tenantID == "" {
		return appErr.ErrTenantIDRequired
	}
	if branchID == "" {
		return appErr.ErrBranchIDRequired
	}
	if productID == "" {
		return appErr.ErrProductIDRequired
	}
	return nil
}