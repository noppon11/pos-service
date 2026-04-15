package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"pos-service/internal/domain"
	"pos-service/internal/dto"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPosService struct {
	mock.Mock
}

func (m *MockPosService) GetHealth(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockPosService) GetHealthByTenantID(ctx context.Context, tenantID string) error {
	args := m.Called(ctx, tenantID)
	return args.Error(0)
}

func (m *MockPosService) GetBranchesByTenantID(ctx context.Context, tenantID string) ([]domain.BranchResponse, error) {
	args := m.Called(ctx, tenantID)

	var data []domain.BranchResponse
	if v := args.Get(0); v != nil {
		data = v.([]domain.BranchResponse)
	}

	return data, args.Error(1)
}

func (m *MockPosService) GetBranchDetail(ctx context.Context, tenantID, branchID string) (*domain.BranchResponse, error) {
	args := m.Called(ctx, tenantID, branchID)

	var data *domain.BranchResponse
	if v := args.Get(0); v != nil {
		data = v.(*domain.BranchResponse)
	}

	return data, args.Error(1)
}

func (m *MockPosService) GetProducts(ctx context.Context, tenantID, branchID string) ([]domain.ProductResponse, error) {
	args := m.Called(ctx, tenantID, branchID)

	var data []domain.ProductResponse
	if v := args.Get(0); v != nil {
		data = v.([]domain.ProductResponse)
	}

	return data, args.Error(1)
}

func (m *MockPosService) GetProductByID(ctx context.Context, tenantID, branchID, productID string) (*domain.ProductResponse, error) {
	args := m.Called(ctx, tenantID, branchID, productID)

	var data *domain.ProductResponse
	if v := args.Get(0); v != nil {
		data = v.(*domain.ProductResponse)
	}

	return data, args.Error(1)
}

func (m *MockPosService) CreateNewProduct(ctx context.Context, tenantID, branchID string, req dto.CreateProductRequest) (*domain.ProductResponse, error) {
	args := m.Called(ctx, tenantID, branchID, req)

	var data *domain.ProductResponse
	if v := args.Get(0); v != nil {
		data = v.(*domain.ProductResponse)
	}

	return data, args.Error(1)
}

func (m *MockPosService) UpdateProduct(ctx context.Context, tenantID, branchID, productID string, req dto.UpdateProductRequest) (*domain.ProductResponse, error) {
	args := m.Called(ctx, tenantID, branchID, productID, req)

	var data *domain.ProductResponse
	if v := args.Get(0); v != nil {
		data = v.(*domain.ProductResponse)
	}

	return data, args.Error(1)
}

func (m *MockPosService) DeleteProduct(ctx context.Context, tenantID, branchID, productID string) error {
	args := m.Called(ctx, tenantID, branchID, productID)
	return args.Error(0)
}

type MockValidator struct {
	mock.Mock
}

func (m *MockValidator) TenantIDValidation(tenantID string) error {
	args := m.Called(tenantID)
	return args.Error(0)
}

func (m *MockValidator) BranchIDValidation(branchID string) error {
	args := m.Called(branchID)
	return args.Error(0)
}

func (m *MockValidator) ProductIDValidation(productID string) error {
	args := m.Called(productID)
	return args.Error(0)
}

func setupGinContext(method, path string, body []byte) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	var req *http.Request
	if body != nil {
		req = httptest.NewRequest(method, path, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	c.Request = req

	return c, w
}

func TestGetHealth_Success(t *testing.T) {
	mockService := new(MockPosService)
	mockValidator := new(MockValidator)

	mockService.On("GetHealth", mock.Anything).Return(nil).Once()

	h := NewPosHandler(mockService, mockValidator)

	c, w := setupGinContext(http.MethodGet, "/health", nil)
	h.GetHealth(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "pos-service", resp["service"])
	assert.Equal(t, "ok", resp["status"])
	assert.NotNil(t, resp["timestamp"])

	mockService.AssertExpectations(t)
}

func TestGetHealth_ServiceUnavailable(t *testing.T) {
	mockService := new(MockPosService)
	mockValidator := new(MockValidator)

	mockService.On("GetHealth", mock.Anything).Return(errors.New("db down")).Once()

	h := NewPosHandler(mockService, mockValidator)

	c, w := setupGinContext(http.MethodGet, "/health", nil)
	h.GetHealth(c)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "pos-service", resp["service"])
	assert.Equal(t, "unhealthy", resp["status"])
	assert.Equal(t, "db down", resp["error"])
	assert.NotNil(t, resp["timestamp"])

	mockService.AssertExpectations(t)
}

func TestGetHealthByTenantID_Success(t *testing.T) {
	mockService := new(MockPosService)
	mockValidator := new(MockValidator)

	tenantID := "tenant_001"

	mockValidator.On("TenantIDValidation", tenantID).Return(nil).Once()
	mockService.On("GetHealthByTenantID", mock.Anything, tenantID).Return(nil).Once()

	h := NewPosHandler(mockService, mockValidator)

	c, w := setupGinContext(http.MethodGet, "/api/v1/tenants/"+tenantID+"/health", nil)
	c.Params = gin.Params{{Key: "tenant_id", Value: tenantID}}

	h.GetHealthByTenantID(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "pos-service", resp["service"])
	assert.Equal(t, "ok", resp["status"])
	assert.Equal(t, tenantID, resp["tenant_id"])
	assert.NotNil(t, resp["timestamp"])

	mockValidator.AssertExpectations(t)
	mockService.AssertExpectations(t)
}

func TestGetBranchesByTenantID_Success(t *testing.T) {
	mockService := new(MockPosService)
	mockValidator := new(MockValidator)

	tenantID := "aura-bkk"
	branches := []domain.BranchResponse{
		{
			BranchID:   "bkk-001",
			BranchName: "Aura Siam",
			Status:     "active",
			Timezone:   "Asia/Bangkok",
			Currency:   "THB",
		},
		{
			BranchID:   "bkk-002",
			BranchName: "Aura Ari",
			Status:     "inactive",
			Timezone:   "Asia/Bangkok",
			Currency:   "THB",
		},
	}

	mockValidator.On("TenantIDValidation", tenantID).Return(nil).Once()
	mockService.On("GetBranchesByTenantID", mock.Anything, tenantID).Return(branches, nil).Once()

	h := NewPosHandler(mockService, mockValidator)

	c, w := setupGinContext(http.MethodGet, "/api/v1/tenants/"+tenantID+"/branches", nil)
	c.Params = gin.Params{{Key: "tenant_id", Value: tenantID}}

	h.GetBranchesByTenantID(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockValidator.AssertExpectations(t)
	mockService.AssertExpectations(t)
}

func TestGetByTenantIDAndBranchID_Success(t *testing.T) {
	mockService := new(MockPosService)
	mockValidator := new(MockValidator)

	tenantID := "aura-bkk"
	branchID := "bkk-001"
	branch := &domain.BranchResponse{
		BranchID:   branchID,
		BranchName: "Aura Siam",
		Status:     "active",
		Timezone:   "Asia/Bangkok",
		Currency:   "THB",
	}

	mockValidator.On("TenantIDValidation", tenantID).Return(nil).Once()
	mockValidator.On("BranchIDValidation", branchID).Return(nil).Once()
	mockService.On("GetBranchDetail", mock.Anything, tenantID, branchID).Return(branch, nil).Once()

	h := NewPosHandler(mockService, mockValidator)

	c, w := setupGinContext(http.MethodGet, "/api/v1/tenants/"+tenantID+"/branches/"+branchID, nil)
	c.Params = gin.Params{
		{Key: "tenant_id", Value: tenantID},
		{Key: "branch_id", Value: branchID},
	}

	h.GetByTenantIDAndBranchID(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp domain.BranchResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, branchID, resp.BranchID)
	assert.Equal(t, "Aura Siam", resp.BranchName)

	mockValidator.AssertExpectations(t)
	mockService.AssertExpectations(t)
}

func TestGetAllProducts_Success(t *testing.T) {
	mockService := new(MockPosService)
	mockValidator := new(MockValidator)

	tenantID := "aura-bkk"
	branchID := "bkk-001"

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

	mockValidator.On("TenantIDValidation", tenantID).Return(nil).Once()
	mockValidator.On("BranchIDValidation", branchID).Return(nil).Once()
	mockService.On("GetProducts", mock.Anything, tenantID, branchID).Return(products, nil).Once()

	h := NewPosHandler(mockService, mockValidator)

	c, w := setupGinContext(http.MethodGet, "/api/v1/tenants/"+tenantID+"/branches/"+branchID+"/products", nil)
	c.Params = gin.Params{
		{Key: "tenant_id", Value: tenantID},
		{Key: "branch_id", Value: branchID},
	}

	h.GetAllProducts(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp []domain.ProductResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Len(t, resp, 1)
	assert.Equal(t, "prod-001", resp[0].ProductID)

	mockValidator.AssertExpectations(t)
	mockService.AssertExpectations(t)
}

func TestGetProductByID_Success(t *testing.T) {
	mockService := new(MockPosService)
	mockValidator := new(MockValidator)

	tenantID := "aura-bkk"
	branchID := "bkk-001"
	productID := "prod-001"

	product := &domain.ProductResponse{
		ProductID:  productID,
		Name:       "Botox 50u",
		SKU:        "BOT-50",
		Price:      3500,
		CategoryID: "treatment",
		Unit:       "unit",
		IsActive:   true,
	}

	mockValidator.On("TenantIDValidation", tenantID).Return(nil).Once()
	mockValidator.On("BranchIDValidation", branchID).Return(nil).Once()
	mockValidator.On("ProductIDValidation", productID).Return(nil).Once()
	mockService.On("GetProductByID", mock.Anything, tenantID, branchID, productID).Return(product, nil).Once()

	h := NewPosHandler(mockService, mockValidator)

	c, w := setupGinContext(http.MethodGet, "/api/v1/tenants/"+tenantID+"/branches/"+branchID+"/products/"+productID, nil)
	c.Params = gin.Params{
		{Key: "tenant_id", Value: tenantID},
		{Key: "branch_id", Value: branchID},
		{Key: "product_id", Value: productID},
	}

	h.GetProductByID(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp domain.ProductResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, productID, resp.ProductID)

	mockValidator.AssertExpectations(t)
	mockService.AssertExpectations(t)
}

func TestCreateProduct_Success(t *testing.T) {
	mockService := new(MockPosService)
	mockValidator := new(MockValidator)

	tenantID := "aura-bkk"
	branchID := "bkk-001"

	req := dto.CreateProductRequest{
		Name:       "Botox 100u",
		SKU:        "BOT-100",
		Price:      6500,
		CategoryID: "treatment",
		Unit:       "unit",
		IsActive:   true,
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

	mockValidator.On("TenantIDValidation", tenantID).Return(nil).Once()
	mockValidator.On("BranchIDValidation", branchID).Return(nil).Once()
	mockService.On("CreateNewProduct", mock.Anything, tenantID, branchID, req).Return(created, nil).Once()

	h := NewPosHandler(mockService, mockValidator)

	body, _ := json.Marshal(req)
	c, w := setupGinContext(http.MethodPost, "/api/v1/tenants/"+tenantID+"/branches/"+branchID+"/products", body)
	c.Params = gin.Params{
		{Key: "tenant_id", Value: tenantID},
		{Key: "branch_id", Value: branchID},
	}

	h.CreateProduct(c)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp domain.ProductResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "prod-001", resp.ProductID)
	assert.Equal(t, "Botox 100u", resp.Name)

	mockValidator.AssertExpectations(t)
	mockService.AssertExpectations(t)
}

func TestCreateProduct_BindJSONError(t *testing.T) {
	mockService := new(MockPosService)
	mockValidator := new(MockValidator)

	tenantID := "aura-bkk"
	branchID := "bkk-001"

	mockValidator.On("TenantIDValidation", tenantID).Return(nil).Once()
	mockValidator.On("BranchIDValidation", branchID).Return(nil).Once()

	h := NewPosHandler(mockService, mockValidator)

	body := []byte(`{"name":`)
	c, w := setupGinContext(http.MethodPost, "/api/v1/tenants/"+tenantID+"/branches/"+branchID+"/products", body)
	c.Params = gin.Params{
		{Key: "tenant_id", Value: tenantID},
		{Key: "branch_id", Value: branchID},
	}

	h.CreateProduct(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	mockValidator.AssertExpectations(t)
	mockService.AssertNotCalled(t, "CreateNewProduct", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestUpdateProduct_Success(t *testing.T) {
	mockService := new(MockPosService)
	mockValidator := new(MockValidator)

	tenantID := "aura-bkk"
	branchID := "bkk-001"
	productID := "86a915ed-d62f-4d91-bba4-372f55b4bbd4"

	req := dto.UpdateProductRequest{
		ProductID:  "86a915ed-d62f-4d91-bba4-372f55b4bbd4",
		Name:       "Botox 120u",
		SKU:        "BOT-120",
		Price:      7500,
		CategoryID: "treatment",
		Unit:       "unit",
		IsActive:   true,
	}

	updated := &domain.ProductResponse{
		ProductID:  productID,
		Name:       req.Name,
		SKU:        req.SKU,
		Price:      req.Price,
		CategoryID: req.CategoryID,
		Unit:       req.Unit,
		IsActive:   req.IsActive,
	}

	mockValidator.On("TenantIDValidation", tenantID).Return(nil).Once()
	mockValidator.On("BranchIDValidation", branchID).Return(nil).Once()
	mockValidator.On("ProductIDValidation", productID).Return(nil).Once()
	mockService.On("UpdateProduct", mock.Anything, tenantID, branchID, productID, req).Return(updated, nil).Once()

	h := NewPosHandler(mockService, mockValidator)

	body, _ := json.Marshal(req)
	c, w := setupGinContext(http.MethodPut, "/api/v1/tenants/"+tenantID+"/branches/"+branchID+"/products/"+productID, body)
	c.Params = gin.Params{
		{Key: "tenant_id", Value: tenantID},
		{Key: "branch_id", Value: branchID},
		{Key: "product_id", Value: productID},
	}

	h.UpdateProduct(c)

	assert.Equal(t, http.StatusOK, w.Code, w.Body.String())

	var resp domain.ProductResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, productID, resp.ProductID)
	assert.Equal(t, "Botox 120u", resp.Name)

	mockValidator.AssertExpectations(t)
	mockService.AssertExpectations(t)
}

func TestDeleteProduct_Success(t *testing.T) {
	mockService := new(MockPosService)
	mockValidator := new(MockValidator)

	tenantID := "aura-bkk"
	branchID := "bkk-001"
	productID := "prod-001"

	mockValidator.On("TenantIDValidation", tenantID).Return(nil).Once()
	mockValidator.On("BranchIDValidation", branchID).Return(nil).Once()
	mockValidator.On("ProductIDValidation", productID).Return(nil).Once()
	mockService.On("DeleteProduct", mock.Anything, tenantID, branchID, productID).Return(nil).Once()

	h := NewPosHandler(mockService, mockValidator)

	c, w := setupGinContext(http.MethodDelete, "/api/v1/tenants/"+tenantID+"/branches/"+branchID+"/products/"+productID, nil)
	c.Params = gin.Params{
		{Key: "tenant_id", Value: tenantID},
		{Key: "branch_id", Value: branchID},
		{Key: "product_id", Value: productID},
	}

	h.DeleteProduct(c)

	assert.True(t, w.Code == http.StatusNoContent || w.Code == http.StatusOK)

	mockValidator.AssertExpectations(t)
	mockService.AssertExpectations(t)
}

func TestToProductResponse_IncludeDeletedAt(t *testing.T) {
	now := time.Now()

	input := &domain.ProductResponse{
		ProductID:  "prod-001",
		Name:       "Botox 50u",
		SKU:        "BOT-50",
		Price:      3500,
		CategoryID: "treatment",
		Unit:       "unit",
		IsActive:   true,
		DeletedAt:  &now,
	}

	out := toProductResponse(input)
	assert.NotNil(t, out)
	assert.Equal(t, "prod-001", out.ProductID)
	assert.NotNil(t, out.DeletedAt)
}

func TestGetAllProducts_InvalidBranchID(t *testing.T) {
	mockService := new(MockPosService)
	mockValidator := new(MockValidator)

	tenantID := "aura-bkk"
	branchID := "INVALID_BRANCH"

	mockValidator.On("TenantIDValidation", tenantID).Return(nil).Once()
	mockValidator.On("BranchIDValidation", branchID).
		Return(errors.New("branch_id is invalid")).Once()

	h := NewPosHandler(mockService, mockValidator)

	c, w := setupGinContext(http.MethodGet, "/api/v1/tenants/"+tenantID+"/branches/"+branchID+"/products", nil)
	c.Params = gin.Params{
		{Key: "tenant_id", Value: tenantID},
		{Key: "branch_id", Value: branchID},
	}

	h.GetAllProducts(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "branch_id is invalid", resp["error"])

	mockValidator.AssertExpectations(t)
	mockService.AssertNotCalled(t, "GetProducts", mock.Anything, mock.Anything, mock.Anything)
}

func TestGetAllProducts_InternalError(t *testing.T) {
	mockService := new(MockPosService)
	mockValidator := new(MockValidator)

	tenantID := "aura-bkk"
	branchID := "bkk-001"

	mockValidator.On("BranchIDValidation", branchID).Return(nil).Once()
	mockValidator.On("TenantIDValidation", tenantID).Return(nil).Once()
	mockService.On("GetProducts", mock.Anything, tenantID, branchID).
		Return(nil, errors.New("repository error")).Once()

	h := NewPosHandler(mockService, mockValidator)

	c, w := setupGinContext(http.MethodGet, "/api/v1/tenants/"+tenantID+"/branches/"+branchID+"/products", nil)
	c.Params = gin.Params{
		{Key: "tenant_id", Value: tenantID},
		{Key: "branch_id", Value: branchID},
	}

	h.GetAllProducts(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "repository error", resp["error"])

	mockValidator.AssertExpectations(t)
	mockService.AssertExpectations(t)
}

func TestGetProductByID_InvalidProductID(t *testing.T) {
	mockService := new(MockPosService)
	mockValidator := new(MockValidator)

	tenantID := "aura-bkk"
	branchID := "bkk-001"
	productID := "INVALID_PRODUCT"

	mockValidator.On("BranchIDValidation", branchID).Return(nil).Once()
	mockValidator.On("TenantIDValidation", tenantID).Return(nil).Once()
	mockValidator.On("ProductIDValidation", productID).
		Return(errors.New("product_id is invalid")).Once()

	h := NewPosHandler(mockService, mockValidator)

	c, w := setupGinContext(http.MethodGet, "/api/v1/tenants/"+tenantID+"/branches/"+branchID+"/products/"+productID, nil)
	c.Params = gin.Params{
		{Key: "tenant_id", Value: tenantID},
		{Key: "branch_id", Value: branchID},
		{Key: "product_id", Value: productID},
	}

	h.GetProductByID(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "product_id is invalid", resp["error"])

	mockValidator.AssertExpectations(t)
	mockService.AssertNotCalled(t, "GetProductByID", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestGetProductByID_InternalError(t *testing.T) {
	mockService := new(MockPosService)
	mockValidator := new(MockValidator)

	tenantID := "aura-bkk"
	branchID := "bkk-001"
	productID := "86a915ed-d62f-4d91-bba4-372f55b4bbd4"

	mockValidator.On("BranchIDValidation", branchID).Return(nil).Once()
	mockValidator.On("TenantIDValidation", tenantID).Return(nil).Once()
	mockValidator.On("ProductIDValidation", productID).Return(nil).Once()
	mockService.On("GetProductByID", mock.Anything, tenantID, branchID, productID).
		Return(nil, errors.New("repository error")).Once()

	h := NewPosHandler(mockService, mockValidator)

	c, w := setupGinContext(http.MethodGet, "/api/v1/tenants/"+tenantID+"/branches/"+branchID+"/products/"+productID, nil)
	c.Params = gin.Params{
		{Key: "tenant_id", Value: tenantID},
		{Key: "branch_id", Value: branchID},
		{Key: "product_id", Value: productID},
	}

	h.GetProductByID(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "repository error", resp["error"])

	mockValidator.AssertExpectations(t)
	mockService.AssertExpectations(t)
}

func TestCreateProduct_InternalError(t *testing.T) {
	mockService := new(MockPosService)
	mockValidator := new(MockValidator)

	tenantID := "aura-bkk"
	branchID := "bkk-001"

	req := dto.CreateProductRequest{
		Name:       "Botox 100u",
		SKU:        "BOT-100",
		Price:      6500,
		CategoryID: "treatment",
		Unit:       "unit",
		IsActive:   true,
	}

	mockValidator.On("TenantIDValidation", tenantID).Return(nil).Once()
	mockValidator.On("BranchIDValidation", branchID).Return(nil).Once()
	mockService.On("CreateNewProduct", mock.Anything, tenantID, branchID, req).
		Return(nil, errors.New("create failed")).Once()

	h := NewPosHandler(mockService, mockValidator)

	body, _ := json.Marshal(req)
	c, w := setupGinContext(http.MethodPost, "/api/v1/tenants/"+tenantID+"/branches/"+branchID+"/products", body)
	c.Params = gin.Params{
		{Key: "tenant_id", Value: tenantID},
		{Key: "branch_id", Value: branchID},
	}

	h.CreateProduct(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "create failed", resp["error"])

	mockValidator.AssertExpectations(t)
	mockService.AssertExpectations(t)
}

func TestUpdateProduct_BindJSONError(t *testing.T) {
	mockService := new(MockPosService)
	mockValidator := new(MockValidator)

	tenantID := "aura-bkk"
	branchID := "bkk-001"
	productID := "86a915ed-d62f-4d91-bba4-372f55b4bbd4"

	mockValidator.On("TenantIDValidation", tenantID).Return(nil).Once()
	mockValidator.On("BranchIDValidation", branchID).Return(nil).Once()
	mockValidator.On("ProductIDValidation", productID).Return(nil).Once()

	h := NewPosHandler(mockService, mockValidator)

	body := []byte(`{"name":`)
	c, w := setupGinContext(http.MethodPut, "/api/v1/tenants/"+tenantID+"/branches/"+branchID+"/products/"+productID, body)
	c.Params = gin.Params{
		{Key: "tenant_id", Value: tenantID},
		{Key: "branch_id", Value: branchID},
		{Key: "product_id", Value: productID},
	}

	h.UpdateProduct(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	mockValidator.AssertExpectations(t)
	mockService.AssertNotCalled(t, "UpdateProduct", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestUpdateProduct_InternalError(t *testing.T) {
	mockService := new(MockPosService)
	mockValidator := new(MockValidator)

	tenantID := "aura-bkk"
	branchID := "bkk-001"
	productID := "86a915ed-d62f-4d91-bba4-372f55b4bbd4"

	req := dto.UpdateProductRequest{
		ProductID:  productID,
		Name:       "Botox 120u",
		SKU:        "BOT-120",
		Price:      7500,
		CategoryID: "treatment",
		Unit:       "unit",
		IsActive:   true,
	}

	mockValidator.On("TenantIDValidation", tenantID).Return(nil).Once()
	mockValidator.On("BranchIDValidation", branchID).Return(nil).Once()
	mockValidator.On("ProductIDValidation", productID).Return(nil).Once()
	mockService.On("UpdateProduct", mock.Anything, tenantID, branchID, productID, req).
		Return(nil, errors.New("update failed")).Once()

	h := NewPosHandler(mockService, mockValidator)

	body, _ := json.Marshal(req)
	c, w := setupGinContext(http.MethodPut, "/api/v1/tenants/"+tenantID+"/branches/"+branchID+"/products/"+productID, body)
	c.Params = gin.Params{
		{Key: "tenant_id", Value: tenantID},
		{Key: "branch_id", Value: branchID},
		{Key: "product_id", Value: productID},
	}

	h.UpdateProduct(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "update failed", resp["error"])

	mockValidator.AssertExpectations(t)
	mockService.AssertExpectations(t)
}

func TestDeleteProduct_InternalError(t *testing.T) {
	mockService := new(MockPosService)
	mockValidator := new(MockValidator)

	tenantID := "aura-bkk"
	branchID := "bkk-001"
	productID := "86a915ed-d62f-4d91-bba4-372f55b4bbd4"

	mockValidator.On("TenantIDValidation", tenantID).Return(nil).Once()
	mockValidator.On("BranchIDValidation", branchID).Return(nil).Once()
	mockValidator.On("ProductIDValidation", productID).Return(nil).Once()
	mockService.On("DeleteProduct", mock.Anything, tenantID, branchID, productID).
		Return(errors.New("delete failed")).Once()

	h := NewPosHandler(mockService, mockValidator)

	c, w := setupGinContext(http.MethodDelete, "/api/v1/tenants/"+tenantID+"/branches/"+branchID+"/products/"+productID, nil)
	c.Params = gin.Params{
		{Key: "tenant_id", Value: tenantID},
		{Key: "branch_id", Value: branchID},
		{Key: "product_id", Value: productID},
	}

	h.DeleteProduct(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "delete failed", resp["error"])

	mockValidator.AssertExpectations(t)
	mockService.AssertExpectations(t)
}

func TestDeleteProduct_InvalidProductID(t *testing.T) {
	mockService := new(MockPosService)
	mockValidator := new(MockValidator)

	tenantID := "aura-bkk"
	branchID := "bkk-001"
	productID := "INVALID_PRODUCT"

	mockValidator.On("TenantIDValidation", tenantID).Return(nil).Once()
	mockValidator.On("BranchIDValidation", branchID).Return(nil).Once()
	mockValidator.On("ProductIDValidation", productID).
		Return(errors.New("product_id is invalid")).Once()

	h := NewPosHandler(mockService, mockValidator)

	c, w := setupGinContext(http.MethodDelete, "/api/v1/tenants/"+tenantID+"/branches/"+branchID+"/products/"+productID, nil)
	c.Params = gin.Params{
		{Key: "tenant_id", Value: tenantID},
		{Key: "branch_id", Value: branchID},
		{Key: "product_id", Value: productID},
	}

	h.DeleteProduct(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "product_id is invalid", resp["error"])

	mockValidator.AssertExpectations(t)
	mockService.AssertNotCalled(t, "DeleteProduct", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestReadiness_Success(t *testing.T) {
	mockService := new(MockPosService)
	mockValidator := new(MockValidator)

	mockService.On("GetHealth", mock.Anything).Return(nil).Once()

	h := NewPosHandler(mockService, mockValidator)

	c, w := setupGinContext(http.MethodGet, "/readiness", nil)
	h.Readiness(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "ready", resp["status"])

	mockService.AssertExpectations(t)
}

func TestReadiness_NotReady(t *testing.T) {
	mockService := new(MockPosService)
	mockValidator := new(MockValidator)

	mockService.On("GetHealth", mock.Anything).Return(errors.New("db down")).Once()

	h := NewPosHandler(mockService, mockValidator)

	c, w := setupGinContext(http.MethodGet, "/readiness", nil)
	h.Readiness(c)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)

	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "not_ready", resp["status"])
	assert.Equal(t, "db down", resp["error"])

	mockService.AssertExpectations(t)
}

func TestGetAllProducts_InvalidTenantID(t *testing.T) {
	mockService := new(MockPosService)
	mockValidator := new(MockValidator)

	tenantID := "AURA"
	branchID := "bkk-001"

	mockValidator.On("TenantIDValidation", tenantID).
		Return(errors.New("tenant_id must be 3-50 chars, lowercase letters, numbers, underscore or dash only")).
		Once()

	h := NewPosHandler(mockService, mockValidator)

	c, w := setupGinContext(http.MethodGet, "/api/v1/tenants/"+tenantID+"/branches/"+branchID+"/products", nil)
	c.Params = gin.Params{
		{Key: "tenant_id", Value: tenantID},
		{Key: "branch_id", Value: branchID},
	}

	h.GetAllProducts(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Contains(t, resp["error"], "tenant_id")

	mockValidator.AssertExpectations(t)
	mockService.AssertNotCalled(t, "GetProducts", mock.Anything, mock.Anything, mock.Anything)
}