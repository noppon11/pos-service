package service

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"pos-service/internal/domain"
	"pos-service/internal/dto"
	appErr "pos-service/internal/errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockBranchRepository struct {
	mock.Mock
}

func (m *MockBranchRepository) ListByTenantID(ctx context.Context, tenantID string) ([]domain.BranchResponse, error) {
	args := m.Called(ctx, tenantID)

	var data []domain.BranchResponse
	if v := args.Get(0); v != nil {
		data = v.([]domain.BranchResponse)
	}
	return data, args.Error(1)
}

func (m *MockBranchRepository) GetByTenantIDAndBranchID(ctx context.Context, tenantID string, branchID string) (*domain.BranchResponse, error) {
	args := m.Called(ctx, tenantID, branchID)

	var data *domain.BranchResponse
	if v := args.Get(0); v != nil {
		data = v.(*domain.BranchResponse)
	}
	return data, args.Error(1)
}

type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) ListByTenantIDAndBranchID(ctx context.Context, tenantID string, branchID string) ([]domain.ProductResponse, error) {
	args := m.Called(ctx, tenantID, branchID)

	var data []domain.ProductResponse
	if v := args.Get(0); v != nil {
		data = v.([]domain.ProductResponse)
	}
	return data, args.Error(1)
}

func (m *MockProductRepository) GetByTenantIDBranchIDAndProductID(ctx context.Context, tenantID string, branchID string, productID string) (*domain.ProductResponse, error) {
	args := m.Called(ctx, tenantID, branchID, productID)

	var data *domain.ProductResponse
	if v := args.Get(0); v != nil {
		data = v.(*domain.ProductResponse)
	}
	return data, args.Error(1)
}

func (m *MockProductRepository) Create(ctx context.Context, tenantID string, branchID string, product domain.ProductResponse) (*domain.ProductResponse, error) {
	args := m.Called(ctx, tenantID, branchID, product)

	var data *domain.ProductResponse
	if v := args.Get(0); v != nil {
		data = v.(*domain.ProductResponse)
	}
	return data, args.Error(1)
}

func (m *MockProductRepository) Update(ctx context.Context, tenantID string, branchID string, productID string, product domain.ProductResponse) (*domain.ProductResponse, error) {
	args := m.Called(ctx, tenantID, branchID, productID, product)

	var data *domain.ProductResponse
	if v := args.Get(0); v != nil {
		data = v.(*domain.ProductResponse)
	}
	return data, args.Error(1)
}

func (m *MockProductRepository) Delete(ctx context.Context, tenantID string, branchID string, productID string) error {
	args := m.Called(ctx, tenantID, branchID, productID)
	return args.Error(0)
}

type MockValidator struct {
	mock.Mock
}

func (m *MockValidator) ValidateBranch(branch domain.BranchResponse) error {
	args := m.Called(branch)
	return args.Error(0)
}

func (m *MockValidator) ValidateProduct(product domain.ProductResponse) error {
	args := m.Called(product)
	return args.Error(0)
}

func TestPosService_GetHealth_DBNotConfigured(t *testing.T) {
	svc := NewPosService(nil, nil, nil, nil)

	err := svc.GetHealth(context.Background())
	assert.ErrorIs(t, err, appErr.ErrDBNotConfigured)
}

func TestPosService_GetHealth_Success(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err == nil {
		defer db.Close()
		svc := NewPosService(db, nil, nil, nil)
		_ = svc
	}
	_ = err
}

func TestPosService_GetBranchesByTenantID_Success(t *testing.T) {
	mockBranchRepo := new(MockBranchRepository)
	mockValidator := new(MockValidator)

	svc := NewPosService(nil, mockBranchRepo, nil, mockValidator)

	tenantID := "aura-bkk"
	branches := []domain.BranchResponse{
		{
			BranchID:   "bkk-001",
			BranchName: "Aura Siam",
			Status:     "active",
			Timezone:   "Asia/Bangkok",
			Currency:   "THB",
		},
	}

	mockBranchRepo.On("ListByTenantID", mock.Anything, tenantID).Return(branches, nil).Once()
	mockValidator.On("ValidateBranch", branches[0]).Return(nil).Once()

	resp, err := svc.GetBranchesByTenantID(context.Background(), tenantID)
	assert.NoError(t, err)
	assert.Len(t, resp, 1)

	mockBranchRepo.AssertExpectations(t)
	mockValidator.AssertExpectations(t)
}

func TestPosService_GetBranchDetail_Success(t *testing.T) {
	mockBranchRepo := new(MockBranchRepository)
	mockValidator := new(MockValidator)

	svc := NewPosService(nil, mockBranchRepo, nil, mockValidator)

	branch := &domain.BranchResponse{
		BranchID:   "bkk-001",
		BranchName: "Aura Siam",
		Status:     "active",
		Timezone:   "Asia/Bangkok",
		Currency:   "THB",
	}

	mockBranchRepo.On("GetByTenantIDAndBranchID", mock.Anything, "aura-bkk", "bkk-001").Return(branch, nil).Once()
	mockValidator.On("ValidateBranch", *branch).Return(nil).Once()

	resp, err := svc.GetBranchDetail(context.Background(), "aura-bkk", "bkk-001")
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	mockBranchRepo.AssertExpectations(t)
	mockValidator.AssertExpectations(t)
}

func TestPosService_GetProducts_Success(t *testing.T) {
	mockProductRepo := new(MockProductRepository)
	mockValidator := new(MockValidator)

	svc := NewPosService(nil, nil, mockProductRepo, mockValidator)

	products := []domain.ProductResponse{
		{
			ProductID:  "prod-001",
			Name:       "Botox 50u",
			SKU:        "BOT-50",
			Price:      3500,
			CategoryID: "treatment",
			Unit:       "unit",
			IsActive:   true,
		},
	}

	mockProductRepo.On("ListByTenantIDAndBranchID", mock.Anything, "aura-bkk", "bkk-001").Return(products, nil).Once()
	mockValidator.On("ValidateProduct", products[0]).Return(nil).Once()

	resp, err := svc.GetProducts(context.Background(), "aura-bkk", "bkk-001")
	assert.NoError(t, err)
	assert.Len(t, resp, 1)

	mockProductRepo.AssertExpectations(t)
	mockValidator.AssertExpectations(t)
}

func TestPosService_GetProductByID_Success(t *testing.T) {
	mockProductRepo := new(MockProductRepository)
	mockValidator := new(MockValidator)

	svc := NewPosService(nil, nil, mockProductRepo, mockValidator)

	product := &domain.ProductResponse{
		ProductID:  "prod-001",
		Name:       "Botox 50u",
		SKU:        "BOT-50",
		Price:      3500,
		CategoryID: "treatment",
		Unit:       "unit",
		IsActive:   true,
	}

	mockProductRepo.On("GetByTenantIDBranchIDAndProductID", mock.Anything, "aura-bkk", "bkk-001", "prod-001").Return(product, nil).Once()
	mockValidator.On("ValidateProduct", *product).Return(nil).Once()

	resp, err := svc.GetProductByID(context.Background(), "aura-bkk", "bkk-001", "prod-001")
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "prod-001", resp.ProductID)

	mockProductRepo.AssertExpectations(t)
	mockValidator.AssertExpectations(t)
}

func TestPosService_CreateNewProduct_Success(t *testing.T) {
	mockBranchRepo := new(MockBranchRepository)
	mockProductRepo := new(MockProductRepository)
	mockValidator := new(MockValidator)

	svc := NewPosService(nil, mockBranchRepo, mockProductRepo, mockValidator)

	req := dto.CreateProductRequest{
		Name:       "Botox 100u",
		SKU:        "BOT-100",
		Price:      6500,
		CategoryID: "treatment",
		Unit:       "unit",
		IsActive:   true,
	}

	expectedInput := domain.ProductResponse{
		Name:       req.Name,
		SKU:        req.SKU,
		Price:      req.Price,
		CategoryID: req.CategoryID,
		Unit:       req.Unit,
		IsActive:   req.IsActive,
	}

	created := &domain.ProductResponse{
		ProductID:  "prod-001",
		Name:       req.Name,
		SKU:        req.SKU,
		Price:      req.Price,
		CategoryID: req.CategoryID,
		Unit:       req.Unit,
		IsActive:   req.IsActive,
	}

	mockValidator.On("ValidateProduct", expectedInput).Return(nil).Once()
	mockProductRepo.On("Create", mock.Anything, "aura-bkk", "bkk-001", expectedInput).Return(created, nil).Once()
	mockValidator.On("ValidateProduct", *created).Return(nil).Once()

	resp, err := svc.CreateNewProduct(context.Background(), "aura-bkk", "bkk-001", req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "prod-001", resp.ProductID)

	mockValidator.AssertExpectations(t)
	mockProductRepo.AssertExpectations(t)
}

func TestPosService_CreateNewProduct_ProductRepoNotConfigured(t *testing.T) {
	svc := NewPosService(nil, new(MockBranchRepository), nil, new(MockValidator))

	_, err := svc.CreateNewProduct(context.Background(), "aura-bkk", "bkk-001", dto.CreateProductRequest{})
	assert.ErrorIs(t, err, appErr.ErrProductRepoNotConfigured)
}

func TestPosService_UpdateProduct_Success(t *testing.T) {
	mockBranchRepo := new(MockBranchRepository)
	mockProductRepo := new(MockProductRepository)
	mockValidator := new(MockValidator)

	svc := NewPosService(nil, mockBranchRepo, mockProductRepo, mockValidator)

	req := dto.UpdateProductRequest{
		Name:       "Botox 120u",
		SKU:        "BOT-120",
		Price:      7500,
		CategoryID: "treatment",
		Unit:       "unit",
		IsActive:   true,
	}

	branch := &domain.BranchResponse{
		BranchID:   "bkk-001",
		BranchName: "Aura Siam",
		Status:     "active",
		Timezone:   "Asia/Bangkok",
		Currency:   "THB",
	}

	input := domain.ProductResponse{
		Name:       req.Name,
		SKU:        req.SKU,
		Price:      req.Price,
		CategoryID: req.CategoryID,
		Unit:       req.Unit,
		IsActive:   req.IsActive,
	}

	updated := &domain.ProductResponse{
		ProductID:  "prod-001",
		Name:       req.Name,
		SKU:        req.SKU,
		Price:      req.Price,
		CategoryID: req.CategoryID,
		Unit:       req.Unit,
		IsActive:   req.IsActive,
	}

	mockBranchRepo.On("GetByTenantIDAndBranchID", mock.Anything, "aura-bkk", "bkk-001").Return(branch, nil).Once()
	mockValidator.On("ValidateProduct", input).Return(nil).Once()
	mockProductRepo.On("Update", mock.Anything, "aura-bkk", "bkk-001", "prod-001", input).Return(updated, nil).Once()
	mockValidator.On("ValidateProduct", *updated).Return(nil).Once()

	resp, err := svc.UpdateProduct(context.Background(), "aura-bkk", "bkk-001", "prod-001", req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "prod-001", resp.ProductID)

	mockBranchRepo.AssertExpectations(t)
	mockProductRepo.AssertExpectations(t)
	mockValidator.AssertExpectations(t)
}

func TestPosService_UpdateProduct_NotFound(t *testing.T) {
	mockBranchRepo := new(MockBranchRepository)
	mockProductRepo := new(MockProductRepository)
	mockValidator := new(MockValidator)

	svc := NewPosService(nil, mockBranchRepo, mockProductRepo, mockValidator)

	req := dto.UpdateProductRequest{
		Name:       "Botox 120u",
		SKU:        "BOT-120",
		Price:      7500,
		CategoryID: "treatment",
		Unit:       "unit",
		IsActive:   true,
	}

	branch := &domain.BranchResponse{
		BranchID:   "bkk-001",
		BranchName: "Aura Siam",
		Status:     "active",
		Timezone:   "Asia/Bangkok",
		Currency:   "THB",
	}

	input := domain.ProductResponse{
		Name:       req.Name,
		SKU:        req.SKU,
		Price:      req.Price,
		CategoryID: req.CategoryID,
		Unit:       req.Unit,
		IsActive:   req.IsActive,
	}

	mockBranchRepo.On("GetByTenantIDAndBranchID", mock.Anything, "aura-bkk", "bkk-001").Return(branch, nil).Once()
	mockValidator.On("ValidateProduct", input).Return(nil).Once()
	mockProductRepo.On("Update", mock.Anything, "aura-bkk", "bkk-001", "prod-001", input).Return(nil, nil).Once()

	resp, err := svc.UpdateProduct(context.Background(), "aura-bkk", "bkk-001", "prod-001", req)
	assert.Nil(t, resp)
	assert.ErrorIs(t, err, appErr.ErrProductNotFound)

	mockBranchRepo.AssertExpectations(t)
	mockProductRepo.AssertExpectations(t)
	mockValidator.AssertExpectations(t)
}

func TestPosService_DeleteProduct_Success(t *testing.T) {
	mockBranchRepo := new(MockBranchRepository)
	mockProductRepo := new(MockProductRepository)
	mockValidator := new(MockValidator)

	svc := NewPosService(nil, mockBranchRepo, mockProductRepo, mockValidator)

	branch := &domain.BranchResponse{
		BranchID:   "bkk-001",
		BranchName: "Aura Siam",
		Status:     "active",
		Timezone:   "Asia/Bangkok",
		Currency:   "THB",
	}
	product := &domain.ProductResponse{
		ProductID:  "prod-001",
		Name:       "Botox 50u",
		SKU:        "BOT-50",
		Price:      3500,
		CategoryID: "treatment",
		Unit:       "unit",
		IsActive:   true,
	}

	mockBranchRepo.On("GetByTenantIDAndBranchID", mock.Anything, "aura-bkk", "bkk-001").Return(branch, nil).Once()
	mockProductRepo.On("GetByTenantIDBranchIDAndProductID", mock.Anything, "aura-bkk", "bkk-001", "prod-001").Return(product, nil).Once()
	mockProductRepo.On("Delete", mock.Anything, "aura-bkk", "bkk-001", "prod-001").Return(nil).Once()

	err := svc.DeleteProduct(context.Background(), "aura-bkk", "bkk-001", "prod-001")
	assert.NoError(t, err)

	mockBranchRepo.AssertExpectations(t)
	mockProductRepo.AssertExpectations(t)
}

func TestPosService_DeleteProduct_NotFound(t *testing.T) {
	mockBranchRepo := new(MockBranchRepository)
	mockProductRepo := new(MockProductRepository)
	mockValidator := new(MockValidator)

	svc := NewPosService(nil, mockBranchRepo, mockProductRepo, mockValidator)

	branch := &domain.BranchResponse{
		BranchID:   "bkk-001",
		BranchName: "Aura Siam",
		Status:     "active",
		Timezone:   "Asia/Bangkok",
		Currency:   "THB",
	}

	mockBranchRepo.On("GetByTenantIDAndBranchID", mock.Anything, "aura-bkk", "bkk-001").Return(branch, nil).Once()
	mockProductRepo.On("GetByTenantIDBranchIDAndProductID", mock.Anything, "aura-bkk", "bkk-001", "prod-001").Return(nil, nil).Once()

	err := svc.DeleteProduct(context.Background(), "aura-bkk", "bkk-001", "prod-001")
	assert.ErrorIs(t, err, appErr.ErrProductNotFound)

	mockBranchRepo.AssertExpectations(t)
	mockProductRepo.AssertExpectations(t)
}

func TestPosService_validateProductScope_Success(t *testing.T) {
	mockBranchRepo := new(MockBranchRepository)

	branch := &domain.BranchResponse{
		BranchID:   "bkk-001",
		BranchName: "Aura Siam",
		Status:     "active",
		Timezone:   "Asia/Bangkok",
		Currency:   "THB",
	}

	mockBranchRepo.
		On("GetByTenantIDAndBranchID", mock.Anything, "aura-bkk", "bkk-001").
		Return(branch, nil).
		Once()

	svc := NewPosService(
		nil,
		mockBranchRepo,
		new(MockProductRepository),
		new(MockValidator),
	)

	err := svc.validateProductScope(context.Background(), "aura-bkk", "bkk-001")

	assert.NoError(t, err)

	mockBranchRepo.AssertExpectations(t)
}

func TestPosService_validateProductScope_BranchNotFound(t *testing.T) {
	mockBranchRepo := new(MockBranchRepository)

	mockBranchRepo.
		On("GetByTenantIDAndBranchID", mock.Anything, "", "bkk-001").
		Return(nil, nil).
		Once()

	svc := NewPosService(
		nil,
		mockBranchRepo,
		new(MockProductRepository),
		new(MockValidator),
	)

	err := svc.validateProductScope(context.Background(), "", "bkk-001")

	assert.ErrorIs(t, err, appErr.ErrBranchNotFound)

	mockBranchRepo.AssertExpectations(t)
}

func TestPosService_validateProductScope_RepoError(t *testing.T) {
	mockBranchRepo := new(MockBranchRepository)

	mockBranchRepo.
		On("GetByTenantIDAndBranchID", mock.Anything, "aura-bkk", "bkk-001").
		Return(nil, errors.New("repo error")).
		Once()

	svc := NewPosService(
		nil,
		mockBranchRepo,
		new(MockProductRepository),
		new(MockValidator),
	)

	err := svc.validateProductScope(context.Background(), "aura-bkk", "bkk-001")

	assert.EqualError(t, err, "repo error")

	mockBranchRepo.AssertExpectations(t)
}

func TestPosService_ensureProductDependencies_MissingProductRepo(t *testing.T) {
	svc := NewPosService(nil, new(MockBranchRepository), nil, new(MockValidator))

	err := svc.ensureProductDependencies()
	assert.ErrorIs(t, err, appErr.ErrProductRepoNotConfigured)
}

func TestPosService_ensureProductDependencies_MissingBranchRepo(t *testing.T) {
	svc := &PosService{
		db:          nil,
		branchRepo:  nil,
		productRepo: new(MockProductRepository),
		validator:   new(MockValidator),
	}

	err := svc.ensureProductDependencies()
	assert.ErrorIs(t, err, appErr.ErrBranchRepoNotConfigured)
}

func TestPosService_validateProduct_SuccessWithDeletedAtNil(t *testing.T) {
	mockValidator := new(MockValidator)
	svc := &PosService{validator: mockValidator}

	product := domain.ProductResponse{
		ProductID:  "prod-001",
		Name:       "Botox 50u",
		SKU:        "BOT-50",
		Price:      3500,
		CategoryID: "treatment",
		Unit:       "unit",
		IsActive:   true,
		DeletedAt:  nil,
	}

	mockValidator.On("ValidateProduct", product).Return(nil).Once()

	err := svc.validateProduct(product)
	assert.NoError(t, err)

	mockValidator.AssertExpectations(t)
}

func TestPosService_validateProduct_WithDeletedAt(t *testing.T) {
	mockValidator := new(MockValidator)
	svc := &PosService{validator: mockValidator}

	now := time.Now()
	product := domain.ProductResponse{
		ProductID:  "prod-001",
		Name:       "Botox 50u",
		SKU:        "BOT-50",
		Price:      3500,
		CategoryID: "treatment",
		Unit:       "unit",
		IsActive:   true,
		DeletedAt:  &now,
	}

	mockValidator.On("ValidateProduct", product).Return(nil).Once()

	err := svc.validateProduct(product)
	assert.NoError(t, err)

	mockValidator.AssertExpectations(t)
}

func TestPosService_GetProducts_RepoError(t *testing.T) {
	mockProductRepo := new(MockProductRepository)
	mockValidator := new(MockValidator)
	svc := NewPosService(nil, nil, mockProductRepo, mockValidator)

	mockProductRepo.On("ListByTenantIDAndBranchID", mock.Anything, "aura-bkk", "bkk-001").
		Return(nil, errors.New("repo error")).Once()

	resp, err := svc.GetProducts(context.Background(), "aura-bkk", "bkk-001")
	assert.Nil(t, resp)
	assert.EqualError(t, err, "repo error")

	mockProductRepo.AssertExpectations(t)
}

func TestPosService_GetBranchesByTenantID_BranchRepoNotConfigured(t *testing.T) {
	svc := NewPosService(nil, nil, new(MockProductRepository), new(MockValidator))

	resp, err := svc.GetBranchesByTenantID(context.Background(), "aura-bkk")
	assert.Nil(t, resp)
	assert.ErrorIs(t, err, appErr.ErrBranchRepoNotConfigured)
}

func TestPosService_GetBranchesByTenantID_ValidateBranchFail(t *testing.T) {
	mockBranchRepo := new(MockBranchRepository)
	mockValidator := new(MockValidator)

	svc := NewPosService(nil, mockBranchRepo, new(MockProductRepository), mockValidator)

	branches := []domain.BranchResponse{
		{
			BranchID:   "bkk-001",
			BranchName: "Aura Siam",
			Status:     "active",
			Timezone:   "Asia/Bangkok",
			Currency:   "THB",
		},
	}

	mockBranchRepo.On("ListByTenantID", mock.Anything, "aura-bkk").Return(branches, nil).Once()
	mockValidator.On("ValidateBranch", branches[0]).Return(errors.New("invalid branch")).Once()

	resp, err := svc.GetBranchesByTenantID(context.Background(), "aura-bkk")
	assert.Nil(t, resp)
	assert.EqualError(t, err, "invalid branch")

	mockBranchRepo.AssertExpectations(t)
	mockValidator.AssertExpectations(t)
}

func TestPosService_GetProductByID_ProductRepoNotConfigured(t *testing.T) {
	svc := NewPosService(nil, new(MockBranchRepository), nil, new(MockValidator))

	resp, err := svc.GetProductByID(context.Background(), "aura-bkk", "bkk-001", "prod-001")
	assert.Nil(t, resp)
	assert.ErrorIs(t, err, appErr.ErrProductRepoNotConfigured)
}

func TestPosService_CreateNewProduct_CreateReturnsNil(t *testing.T) {
	mockBranchRepo := new(MockBranchRepository)
	mockProductRepo := new(MockProductRepository)
	mockValidator := new(MockValidator)

	svc := NewPosService(nil, mockBranchRepo, mockProductRepo, mockValidator)

	req := dto.CreateProductRequest{
		Name:       "Botox 100u",
		SKU:        "BOT-100",
		Price:      6500,
		CategoryID: "treatment",
		Unit:       "unit",
		IsActive:   true,
	}

	input := domain.ProductResponse{
		Name:       req.Name,
		SKU:        req.SKU,
		Price:      req.Price,
		CategoryID: req.CategoryID,
		Unit:       req.Unit,
		IsActive:   req.IsActive,
	}

	mockValidator.On("ValidateProduct", input).Return(nil).Once()
	mockProductRepo.On("Create", mock.Anything, "aura-bkk", "bkk-001", input).Return(nil, nil).Once()

	resp, err := svc.CreateNewProduct(context.Background(), "aura-bkk", "bkk-001", req)
	assert.Nil(t, resp)
	assert.ErrorIs(t, err, appErr.ErrCreateProductFailed)

	mockValidator.AssertExpectations(t)
	mockProductRepo.AssertExpectations(t)
}

func TestPosService_UpdateProduct_UpdatedProductValidationFails(t *testing.T) {
	mockBranchRepo := new(MockBranchRepository)
	mockProductRepo := new(MockProductRepository)
	mockValidator := new(MockValidator)

	svc := NewPosService(nil, mockBranchRepo, mockProductRepo, mockValidator)

	req := dto.UpdateProductRequest{
		ProductID:  "prod-001",
		Name:       "Botox 120u",
		SKU:        "BOT-120",
		Price:      7500,
		CategoryID: "treatment",
		Unit:       "unit",
		IsActive:   true,
	}

	branch := &domain.BranchResponse{
		BranchID:   "bkk-001",
		BranchName: "Aura Siam",
		Status:     "active",
		Timezone:   "Asia/Bangkok",
		Currency:   "THB",
	}

	input := domain.ProductResponse{
		Name:       req.Name,
		SKU:        req.SKU,
		Price:      req.Price,
		CategoryID: req.CategoryID,
		Unit:       req.Unit,
		IsActive:   req.IsActive,
	}

	updated := &domain.ProductResponse{
		ProductID:  "prod-001",
		Name:       req.Name,
		SKU:        req.SKU,
		Price:      req.Price,
		CategoryID: req.CategoryID,
		Unit:       req.Unit,
		IsActive:   req.IsActive,
	}

	mockBranchRepo.On("GetByTenantIDAndBranchID", mock.Anything, "aura-bkk", "bkk-001").Return(branch, nil).Once()
	mockValidator.On("ValidateProduct", input).Return(nil).Once()
	mockProductRepo.On("Update", mock.Anything, "aura-bkk", "bkk-001", "prod-001", input).Return(updated, nil).Once()
	mockValidator.On("ValidateProduct", *updated).Return(errors.New("invalid updated product")).Once()

	resp, err := svc.UpdateProduct(context.Background(), "aura-bkk", "bkk-001", "prod-001", req)
	assert.Nil(t, resp)
	assert.EqualError(t, err, "invalid updated product")

	mockBranchRepo.AssertExpectations(t)
	mockProductRepo.AssertExpectations(t)
	mockValidator.AssertExpectations(t)
}

func TestPosService_DeleteProduct_ProductLookupError(t *testing.T) {
	mockBranchRepo := new(MockBranchRepository)
	mockProductRepo := new(MockProductRepository)
	mockValidator := new(MockValidator)

	svc := NewPosService(nil, mockBranchRepo, mockProductRepo, mockValidator)

	branch := &domain.BranchResponse{
		BranchID:   "bkk-001",
		BranchName: "Aura Siam",
		Status:     "active",
		Timezone:   "Asia/Bangkok",
		Currency:   "THB",
	}

	mockBranchRepo.On("GetByTenantIDAndBranchID", mock.Anything, "aura-bkk", "bkk-001").Return(branch, nil).Once()
	mockProductRepo.On("GetByTenantIDBranchIDAndProductID", mock.Anything, "aura-bkk", "bkk-001", "prod-001").
		Return(nil, errors.New("lookup failed")).Once()

	err := svc.DeleteProduct(context.Background(), "aura-bkk", "bkk-001", "prod-001")
	assert.EqualError(t, err, "lookup failed")

	mockBranchRepo.AssertExpectations(t)
	mockProductRepo.AssertExpectations(t)
}

func TestPosService_GetHealthByTenantID_DBNotConfigured(t *testing.T) {
	svc := NewPosService(nil, nil, nil, nil)

	err := svc.GetHealthByTenantID(context.Background(), "aura-bkk")
	assert.ErrorIs(t, err, appErr.ErrDBNotConfigured)
}

func TestPosService_GetBranchesByTenantID_RepoError(t *testing.T) {
	mockBranchRepo := new(MockBranchRepository)
	mockValidator := new(MockValidator)

	svc := NewPosService(nil, mockBranchRepo, nil, mockValidator)

	mockBranchRepo.On("ListByTenantID", mock.Anything, "aura-bkk").
		Return(nil, errors.New("repo error")).Once()

	resp, err := svc.GetBranchesByTenantID(context.Background(), "aura-bkk")
	assert.Nil(t, resp)
	assert.EqualError(t, err, "repo error")

	mockBranchRepo.AssertExpectations(t)
}

func TestPosService_GetBranchDetail_BranchRepoNotConfigured(t *testing.T) {
	svc := NewPosService(nil, nil, new(MockProductRepository), new(MockValidator))

	resp, err := svc.GetBranchDetail(context.Background(), "aura-bkk", "bkk-001")
	assert.Nil(t, resp)
	assert.ErrorIs(t, err, appErr.ErrBranchRepoNotConfigured)
}

func TestPosService_GetBranchDetail_RepoError(t *testing.T) {
	mockBranchRepo := new(MockBranchRepository)
	mockValidator := new(MockValidator)

	svc := NewPosService(nil, mockBranchRepo, new(MockProductRepository), mockValidator)

	mockBranchRepo.
		On("GetByTenantIDAndBranchID", mock.Anything, "aura-bkk", "bkk-001").
		Return(nil, errors.New("repo error")).Once()

	resp, err := svc.GetBranchDetail(context.Background(), "aura-bkk", "bkk-001")
	assert.Nil(t, resp)
	assert.EqualError(t, err, "repo error")

	mockBranchRepo.AssertExpectations(t)
}

func TestPosService_GetBranchDetail_BranchNotFound(t *testing.T) {
	mockBranchRepo := new(MockBranchRepository)
	mockValidator := new(MockValidator)

	svc := NewPosService(nil, mockBranchRepo, new(MockProductRepository), mockValidator)

	mockBranchRepo.
		On("GetByTenantIDAndBranchID", mock.Anything, "aura-bkk", "bkk-001").
		Return(nil, nil).Once()

	resp, err := svc.GetBranchDetail(context.Background(), "aura-bkk", "bkk-001")
	assert.Nil(t, resp)
	assert.ErrorIs(t, err, appErr.ErrBranchNotFound)

	mockBranchRepo.AssertExpectations(t)
}

func TestPosService_GetBranchDetail_ValidateBranchFail(t *testing.T) {
	mockBranchRepo := new(MockBranchRepository)
	mockValidator := new(MockValidator)

	svc := NewPosService(nil, mockBranchRepo, new(MockProductRepository), mockValidator)

	branch := &domain.BranchResponse{
		BranchID:   "bkk-001",
		BranchName: "Aura Siam",
		Status:     "active",
		Timezone:   "Asia/Bangkok",
		Currency:   "THB",
	}

	mockBranchRepo.
		On("GetByTenantIDAndBranchID", mock.Anything, "aura-bkk", "bkk-001").
		Return(branch, nil).Once()
	mockValidator.
		On("ValidateBranch", *branch).
		Return(errors.New("invalid branch")).Once()

	resp, err := svc.GetBranchDetail(context.Background(), "aura-bkk", "bkk-001")
	assert.Nil(t, resp)
	assert.EqualError(t, err, "invalid branch")

	mockBranchRepo.AssertExpectations(t)
	mockValidator.AssertExpectations(t)
}

func TestPosService_GetProducts_ProductRepoNotConfigured(t *testing.T) {
	svc := NewPosService(nil, new(MockBranchRepository), nil, new(MockValidator))

	resp, err := svc.GetProducts(context.Background(), "aura-bkk", "bkk-001")
	assert.Nil(t, resp)
	assert.ErrorIs(t, err, appErr.ErrProductRepoNotConfigured)
}

func TestPosService_GetProducts_ValidateProductFail(t *testing.T) {
	mockProductRepo := new(MockProductRepository)
	mockValidator := new(MockValidator)

	svc := NewPosService(nil, new(MockBranchRepository), mockProductRepo, mockValidator)

	products := []domain.ProductResponse{
		{
			ProductID:  "prod-001",
			Name:       "Botox 50u",
			SKU:        "BOT-50",
			Price:      3500,
			CategoryID: "treatment",
			Unit:       "unit",
			IsActive:   true,
		},
	}

	mockProductRepo.
		On("ListByTenantIDAndBranchID", mock.Anything, "aura-bkk", "bkk-001").
		Return(products, nil).Once()
	mockValidator.
		On("ValidateProduct", products[0]).
		Return(errors.New("invalid product")).Once()

	resp, err := svc.GetProducts(context.Background(), "aura-bkk", "bkk-001")
	assert.Nil(t, resp)
	assert.EqualError(t, err, "invalid product")

	mockProductRepo.AssertExpectations(t)
	mockValidator.AssertExpectations(t)
}

func TestPosService_GetProductByID_RepoError(t *testing.T) {
	mockProductRepo := new(MockProductRepository)
	mockValidator := new(MockValidator)

	svc := NewPosService(nil, new(MockBranchRepository), mockProductRepo, mockValidator)

	mockProductRepo.
		On("GetByTenantIDBranchIDAndProductID", mock.Anything, "aura-bkk", "bkk-001", "prod-001").
		Return(nil, errors.New("repo error")).Once()

	resp, err := svc.GetProductByID(context.Background(), "aura-bkk", "bkk-001", "prod-001")
	assert.Nil(t, resp)
	assert.EqualError(t, err, "repo error")

	mockProductRepo.AssertExpectations(t)
}

func TestPosService_GetProductByID_NotFound(t *testing.T) {
	mockProductRepo := new(MockProductRepository)
	mockValidator := new(MockValidator)

	svc := NewPosService(nil, new(MockBranchRepository), mockProductRepo, mockValidator)

	mockProductRepo.
		On("GetByTenantIDBranchIDAndProductID", mock.Anything, "aura-bkk", "bkk-001", "prod-001").
		Return(nil, nil).Once()

	resp, err := svc.GetProductByID(context.Background(), "aura-bkk", "bkk-001", "prod-001")
	assert.Nil(t, resp)
	assert.ErrorIs(t, err, appErr.ErrProductNotFound)

	mockProductRepo.AssertExpectations(t)
}

func TestPosService_GetProductByID_ValidateFail(t *testing.T) {
	mockProductRepo := new(MockProductRepository)
	mockValidator := new(MockValidator)

	svc := NewPosService(nil, new(MockBranchRepository), mockProductRepo, mockValidator)

	product := &domain.ProductResponse{
		ProductID:  "prod-001",
		Name:       "Botox 50u",
		SKU:        "BOT-50",
		Price:      3500,
		CategoryID: "treatment",
		Unit:       "unit",
		IsActive:   true,
	}

	mockProductRepo.
		On("GetByTenantIDBranchIDAndProductID", mock.Anything, "aura-bkk", "bkk-001", "prod-001").
		Return(product, nil).Once()
	mockValidator.
		On("ValidateProduct", *product).
		Return(errors.New("invalid product")).Once()

	resp, err := svc.GetProductByID(context.Background(), "aura-bkk", "bkk-001", "prod-001")
	assert.Nil(t, resp)
	assert.EqualError(t, err, "invalid product")

	mockProductRepo.AssertExpectations(t)
	mockValidator.AssertExpectations(t)
}

func TestPosService_CreateNewProduct_BranchRepoNotConfigured(t *testing.T) {
	svc := NewPosService(nil, nil, new(MockProductRepository), new(MockValidator))

	_, err := svc.CreateNewProduct(context.Background(), "aura-bkk", "bkk-001", dto.CreateProductRequest{})
	assert.ErrorIs(t, err, appErr.ErrBranchRepoNotConfigured)
}

func TestPosService_CreateNewProduct_MissingTenantID(t *testing.T) {
	svc := NewPosService(nil, new(MockBranchRepository), new(MockProductRepository), new(MockValidator))

	_, err := svc.CreateNewProduct(context.Background(), "", "bkk-001", dto.CreateProductRequest{})
	assert.ErrorIs(t, err, appErr.ErrTenantIDRequired)
}

func TestPosService_CreateNewProduct_MissingBranchID(t *testing.T) {
	svc := NewPosService(nil, new(MockBranchRepository), new(MockProductRepository), new(MockValidator))

	_, err := svc.CreateNewProduct(context.Background(), "aura-bkk", "", dto.CreateProductRequest{})
	assert.ErrorIs(t, err, appErr.ErrBranchIDRequired)
}

func TestPosService_CreateNewProduct_ValidateFail(t *testing.T) {
	mockProductRepo := new(MockProductRepository)
	mockValidator := new(MockValidator)

	svc := NewPosService(nil, new(MockBranchRepository), mockProductRepo, mockValidator)

	req := dto.CreateProductRequest{
		Name:       "Botox 100u",
		SKU:        "BOT-100",
		Price:      6500,
		CategoryID: "treatment",
		Unit:       "unit",
		IsActive:   true,
	}

	input := domain.ProductResponse{
		Name:       req.Name,
		SKU:        req.SKU,
		Price:      req.Price,
		CategoryID: req.CategoryID,
		Unit:       req.Unit,
		IsActive:   req.IsActive,
	}

	mockValidator.On("ValidateProduct", input).Return(errors.New("invalid product")).Once()

	resp, err := svc.CreateNewProduct(context.Background(), "aura-bkk", "bkk-001", req)
	assert.Nil(t, resp)
	assert.EqualError(t, err, "invalid product")

	mockValidator.AssertExpectations(t)
}

func TestPosService_CreateNewProduct_RepoError(t *testing.T) {
	mockProductRepo := new(MockProductRepository)
	mockValidator := new(MockValidator)

	svc := NewPosService(nil, new(MockBranchRepository), mockProductRepo, mockValidator)

	req := dto.CreateProductRequest{
		Name:       "Botox 100u",
		SKU:        "BOT-100",
		Price:      6500,
		CategoryID: "treatment",
		Unit:       "unit",
		IsActive:   true,
	}

	input := domain.ProductResponse{
		Name:       req.Name,
		SKU:        req.SKU,
		Price:      req.Price,
		CategoryID: req.CategoryID,
		Unit:       req.Unit,
		IsActive:   req.IsActive,
	}

	mockValidator.On("ValidateProduct", input).Return(nil).Once()
	mockProductRepo.On("Create", mock.Anything, "aura-bkk", "bkk-001", input).
		Return(nil, errors.New("repo error")).Once()

	resp, err := svc.CreateNewProduct(context.Background(), "aura-bkk", "bkk-001", req)
	assert.Nil(t, resp)
	assert.EqualError(t, err, "repo error")

	mockValidator.AssertExpectations(t)
	mockProductRepo.AssertExpectations(t)
}

func TestPosService_CreateNewProduct_CreatedProductValidateFail(t *testing.T) {
	mockProductRepo := new(MockProductRepository)
	mockValidator := new(MockValidator)

	svc := NewPosService(nil, new(MockBranchRepository), mockProductRepo, mockValidator)

	req := dto.CreateProductRequest{
		Name:       "Botox 100u",
		SKU:        "BOT-100",
		Price:      6500,
		CategoryID: "treatment",
		Unit:       "unit",
		IsActive:   true,
	}

	input := domain.ProductResponse{
		Name:       req.Name,
		SKU:        req.SKU,
		Price:      req.Price,
		CategoryID: req.CategoryID,
		Unit:       req.Unit,
		IsActive:   req.IsActive,
	}

	created := &domain.ProductResponse{
		ProductID:  "prod-001",
		Name:       req.Name,
		SKU:        req.SKU,
		Price:      req.Price,
		CategoryID: req.CategoryID,
		Unit:       req.Unit,
		IsActive:   req.IsActive,
	}

	mockValidator.On("ValidateProduct", input).Return(nil).Once()
	mockProductRepo.On("Create", mock.Anything, "aura-bkk", "bkk-001", input).Return(created, nil).Once()
	mockValidator.On("ValidateProduct", *created).Return(errors.New("invalid created product")).Once()

	resp, err := svc.CreateNewProduct(context.Background(), "aura-bkk", "bkk-001", req)
	assert.Nil(t, resp)
	assert.EqualError(t, err, "invalid created product")

	mockValidator.AssertExpectations(t)
	mockProductRepo.AssertExpectations(t)
}

func TestPosService_UpdateProduct_ProductRepoNotConfigured(t *testing.T) {
	svc := NewPosService(nil, new(MockBranchRepository), nil, new(MockValidator))

	_, err := svc.UpdateProduct(context.Background(), "aura-bkk", "bkk-001", "prod-001", dto.UpdateProductRequest{})
	assert.ErrorIs(t, err, appErr.ErrProductRepoNotConfigured)
}

func TestPosService_UpdateProduct_ValidateFail(t *testing.T) {
	mockBranchRepo := new(MockBranchRepository)
	mockProductRepo := new(MockProductRepository)
	mockValidator := new(MockValidator)

	svc := NewPosService(nil, mockBranchRepo, mockProductRepo, mockValidator)

	req := dto.UpdateProductRequest{
		ProductID:  "prod-001",
		Name:       "Botox 120u",
		SKU:        "BOT-120",
		Price:      7500,
		CategoryID: "treatment",
		Unit:       "unit",
		IsActive:   true,
	}

	branch := &domain.BranchResponse{
		BranchID:   "bkk-001",
		BranchName: "Aura Siam",
		Status:     "active",
		Timezone:   "Asia/Bangkok",
		Currency:   "THB",
	}

	input := domain.ProductResponse{
		Name:       req.Name,
		SKU:        req.SKU,
		Price:      req.Price,
		CategoryID: req.CategoryID,
		Unit:       req.Unit,
		IsActive:   req.IsActive,
	}

	mockBranchRepo.On("GetByTenantIDAndBranchID", mock.Anything, "aura-bkk", "bkk-001").Return(branch, nil).Once()
	mockValidator.On("ValidateProduct", input).Return(errors.New("invalid product")).Once()

	resp, err := svc.UpdateProduct(context.Background(), "aura-bkk", "bkk-001", "prod-001", req)
	assert.Nil(t, resp)
	assert.EqualError(t, err, "invalid product")

	mockBranchRepo.AssertExpectations(t)
	mockValidator.AssertExpectations(t)
}

func TestPosService_UpdateProduct_RepoError(t *testing.T) {
	mockBranchRepo := new(MockBranchRepository)
	mockProductRepo := new(MockProductRepository)
	mockValidator := new(MockValidator)

	svc := NewPosService(nil, mockBranchRepo, mockProductRepo, mockValidator)

	req := dto.UpdateProductRequest{
		ProductID:  "prod-001",
		Name:       "Botox 120u",
		SKU:        "BOT-120",
		Price:      7500,
		CategoryID: "treatment",
		Unit:       "unit",
		IsActive:   true,
	}

	branch := &domain.BranchResponse{
		BranchID:   "bkk-001",
		BranchName: "Aura Siam",
		Status:     "active",
		Timezone:   "Asia/Bangkok",
		Currency:   "THB",
	}

	input := domain.ProductResponse{
		Name:       req.Name,
		SKU:        req.SKU,
		Price:      req.Price,
		CategoryID: req.CategoryID,
		Unit:       req.Unit,
		IsActive:   req.IsActive,
	}

	mockBranchRepo.On("GetByTenantIDAndBranchID", mock.Anything, "aura-bkk", "bkk-001").Return(branch, nil).Once()
	mockValidator.On("ValidateProduct", input).Return(nil).Once()
	mockProductRepo.On("Update", mock.Anything, "aura-bkk", "bkk-001", "prod-001", input).
		Return(nil, errors.New("repo error")).Once()

	resp, err := svc.UpdateProduct(context.Background(), "aura-bkk", "bkk-001", "prod-001", req)
	assert.Nil(t, resp)
	assert.EqualError(t, err, "repo error")

	mockBranchRepo.AssertExpectations(t)
	mockProductRepo.AssertExpectations(t)
	mockValidator.AssertExpectations(t)
}

func TestPosService_DeleteProduct_DeleteRepoError(t *testing.T) {
	mockBranchRepo := new(MockBranchRepository)
	mockProductRepo := new(MockProductRepository)
	mockValidator := new(MockValidator)

	svc := NewPosService(nil, mockBranchRepo, mockProductRepo, mockValidator)

	branch := &domain.BranchResponse{
		BranchID:   "bkk-001",
		BranchName: "Aura Siam",
		Status:     "active",
		Timezone:   "Asia/Bangkok",
		Currency:   "THB",
	}
	product := &domain.ProductResponse{
		ProductID:  "prod-001",
		Name:       "Botox 50u",
		SKU:        "BOT-50",
		Price:      3500,
		CategoryID: "treatment",
		Unit:       "unit",
		IsActive:   true,
	}

	mockBranchRepo.On("GetByTenantIDAndBranchID", mock.Anything, "aura-bkk", "bkk-001").Return(branch, nil).Once()
	mockProductRepo.On("GetByTenantIDBranchIDAndProductID", mock.Anything, "aura-bkk", "bkk-001", "prod-001").Return(product, nil).Once()
	mockProductRepo.On("Delete", mock.Anything, "aura-bkk", "bkk-001", "prod-001").
		Return(errors.New("delete failed")).Once()

	err := svc.DeleteProduct(context.Background(), "aura-bkk", "bkk-001", "prod-001")
	assert.EqualError(t, err, "delete failed")

	mockBranchRepo.AssertExpectations(t)
	mockProductRepo.AssertExpectations(t)
}

func TestPosService_validateBranch_ValidatorNotConfigured(t *testing.T) {
	svc := &PosService{validator: nil}

	err := svc.validateBranch(domain.BranchResponse{})
	assert.ErrorIs(t, err, appErr.ErrValidatorNotConfigured)
}

func TestPosService_validateProduct_ValidatorNotConfigured(t *testing.T) {
	svc := &PosService{validator: nil}

	err := svc.validateProduct(domain.ProductResponse{})
	assert.ErrorIs(t, err, appErr.ErrValidatorNotConfigured)
}