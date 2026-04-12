package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"pos-service/internal/domain"
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

type MockTenantValidator struct {
	mock.Mock
}

func (m *MockTenantValidator) TenantIDValidation(tenantID string) error {
	args := m.Called(tenantID)
	return args.Error(0)
}

func setupGinContext(method, path string) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(method, path, nil)
	c.Request = req
	return c, w
}

func TestGetHealth_Success(t *testing.T) {
	mockService := new(MockPosService)

	mockService.On("GetHealth", mock.Anything).Return(nil).Once()

	h := NewPosHandler(mockService)

	c, w := setupGinContext(http.MethodGet, "/health")
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

	mockService.On("GetHealth", mock.Anything).Return(errors.New("db down")).Once()

	h := NewPosHandler(mockService)

	c, w := setupGinContext(http.MethodGet, "/health")
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

func TestGetHealthByTenantID_MissingTenantID(t *testing.T) {
	mockService := new(MockPosService)

	h := NewPosHandler(mockService)

	c, w := setupGinContext(http.MethodGet, "/tenants//health")
	h.GetHealthByTenantID(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "tenant_id is required", resp["error"])

	mockService.AssertNotCalled(t, "GetHealthByTenantID", mock.Anything, mock.Anything)
}

func TestGetHealthByTenantID_ServiceUnavailable(t *testing.T) {
	mockService := new(MockPosService)
	tenantID := "tenant_001"

	mockService.On("GetHealthByTenantID", mock.Anything, tenantID).Return(errors.New("db down")).Once()

	h := NewPosHandler(mockService)

	c, w := setupGinContext(http.MethodGet, "/tenants/"+tenantID+"/health")
	c.Params = gin.Params{{Key: "tenant_id", Value: tenantID}}

	h.GetHealthByTenantID(c)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "pos-service", resp["service"])
	assert.Equal(t, "unhealthy", resp["status"])
	assert.Equal(t, tenantID, resp["tenant_id"])
	assert.Equal(t, "db down", resp["error"])
	assert.NotNil(t, resp["timestamp"])

	mockService.AssertExpectations(t)
}

func TestGetHealthByTenantID_Success(t *testing.T) {
	mockService := new(MockPosService)
	tenantID := "tenant_001"

	mockService.On("GetHealthByTenantID", mock.Anything, tenantID).Return(nil).Once()

	h := NewPosHandler(mockService)

	c, w := setupGinContext(http.MethodGet, "/tenants/"+tenantID+"/health")
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

	mockService.AssertExpectations(t)
}

func TestGetBranchesByTenantID_Success(t *testing.T) {
	mockService := new(MockPosService)
	tenantID := "aura-bkk"
	branches := []domain.BranchResponse{
		{
			BranchID:   "bkk-001",
			BranchName: "Aura Siam",
			Status:     "active",
		},
		{
			BranchID:   "bkk-002",
			BranchName: "Aura Ari",
			Status:     "inactive",
		},
	}

	mockService.
		On("GetBranchesByTenantID", mock.Anything, tenantID).
		Return(branches, nil).
		Once()

	h := NewPosHandler(mockService)

	c, w := setupGinContext(http.MethodGet, "/api/v1/tenants/"+tenantID+"/branches")
	c.Params = gin.Params{{Key: "tenant_id", Value: tenantID}}

	h.GetBranchesByTenantID(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp domain.ListBranchesResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, tenantID, resp.TenantID)
	assert.Len(t, resp.Data, 2)
	assert.Equal(t, "bkk-001", resp.Data[0].BranchID)
	assert.Equal(t, "Aura Siam", resp.Data[0].BranchName)
	assert.Equal(t, "active", resp.Data[0].Status)

	mockService.AssertExpectations(t)
}

func TestGetBranchesByTenantID_InvalidTenantID(t *testing.T) {
	mockService := new(MockPosService)

	tenantID := "AURA"

	mockService.
		On("GetBranchesByTenantID", mock.Anything, tenantID).
		Return(nil, errors.New("tenant_id must be 3-50 chars, lowercase letters, numbers, underscore or dash only")).
		Once()

	h := NewPosHandler(mockService)

	c, w := setupGinContext(http.MethodGet, "/api/v1/tenants/"+tenantID+"/branches")
	c.Params = gin.Params{{Key: "tenant_id", Value: tenantID}}

	h.GetBranchesByTenantID(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "tenant_id must be 3-50 chars, lowercase letters, numbers, underscore or dash only", resp["error"])

	mockService.AssertExpectations(t)
}

func TestGetBranchesByTenantID_EmptyList(t *testing.T) {
	mockService := new(MockPosService)

	tenantID := "aura-xyz"
	branches := []domain.BranchResponse{}
	mockService.On("GetBranchesByTenantID", mock.Anything, tenantID).Return(branches, nil).Once()

	h := NewPosHandler(mockService)

	c, w := setupGinContext(http.MethodGet, "/api/v1/tenants/"+tenantID+"/branches")
	c.Params = gin.Params{{Key: "tenant_id", Value: tenantID}}

	h.GetBranchesByTenantID(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp domain.ListBranchesResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, tenantID, resp.TenantID)
	assert.Len(t, resp.Data, 0)
	mockService.AssertExpectations(t)
}

func TestGetBranchesByTenantID_InternalError(t *testing.T) {
	mockService := new(MockPosService)
	tenantID := "aura-bkk"

	mockService.On("GetBranchesByTenantID", mock.Anything, tenantID).Return(nil, errors.New("repository error")).Once()

	h := NewPosHandler(mockService)

	c, w := setupGinContext(http.MethodGet, "/api/v1/tenants/"+tenantID+"/branches")
	c.Params = gin.Params{{Key: "tenant_id", Value: tenantID}}

	h.GetBranchesByTenantID(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "repository error", resp["error"])

	mockService.AssertExpectations(t)
}