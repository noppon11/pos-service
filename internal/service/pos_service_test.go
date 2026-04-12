package service

import (
	"context"
	"errors"
	"testing"

	"pos-service/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDB struct {
	mock.Mock
}

type MockRepo struct {
	data []domain.BranchResponse
	err  error
}

func (m *MockRepo) ListByTenantID(ctx context.Context, tenantID string) ([]domain.BranchResponse, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.data, nil
}

func (m *MockDB) PingContext(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

type MockValidator struct {
	mock.Mock
}

func (m *MockValidator) TenantIDValidation(tenantID string) error {
	args := m.Called(tenantID)
	return args.Error(0)
}

func (m *MockValidator) BranchValidation(branch domain.BranchResponse) error {
	args := m.Called(branch)
	return args.Error(0)
}

func TestGetHealth_DBDown(t *testing.T) {
	mockDB := new(MockDB)

	mockDB.
		On("PingContext", mock.Anything).
		Return(errors.New("db down")).
		Once()

	svc := &PosService{db: mockDB}

	err := svc.GetHealth(context.Background())

	assert.Error(t, err)
	assert.EqualError(t, err, "db down")
	mockDB.AssertExpectations(t)
}

func TestGetHealth_Success(t *testing.T) {
	mockDB := new(MockDB)

	mockDB.
		On("PingContext", mock.Anything).
		Return(nil).
		Once()

	svc := &PosService{db: mockDB}

	err := svc.GetHealth(context.Background())

	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
}

func TestGetHealthByTenantID_InvalidTenantID(t *testing.T) {
	mockDB := new(MockDB)
	mockValidator := new(MockValidator)

	mockValidator.
		On("TenantIDValidation", "AURA").
		Return(errors.New("tenant_id must be 3-50 chars, lowercase letters, numbers, underscore or dash only")).
		Once()

	svc := &PosService{
		db:        mockDB,
		validator: mockValidator,
	}

	err := svc.GetHealthByTenantID(context.Background(), "AURA")

	assert.Error(t, err)
	assert.EqualError(t, err, "tenant_id must be 3-50 chars, lowercase letters, numbers, underscore or dash only")

	mockValidator.AssertExpectations(t)
	mockDB.AssertNotCalled(t, "PingContext", mock.Anything)
}

func TestGetHealthByTenantID_DBDown(t *testing.T) {
	mockDB := new(MockDB)
	mockValidator := new(MockValidator)

	mockValidator.
		On("TenantIDValidation", "tenant_001").
		Return(nil).
		Once()

	mockDB.
		On("PingContext", mock.Anything).
		Return(errors.New("db down")).
		Once()

	svc := &PosService{
		db:        mockDB,
		validator: mockValidator,
	}

	err := svc.GetHealthByTenantID(context.Background(), "tenant_001")

	assert.Error(t, err)
	assert.EqualError(t, err, "db down")

	mockValidator.AssertExpectations(t)
	mockDB.AssertExpectations(t)
}

func TestGetHealthByTenantID_Success(t *testing.T) {
	mockDB := new(MockDB)
	mockValidator := new(MockValidator)

	mockValidator.
		On("TenantIDValidation", "tenant_001").
		Return(nil).
		Once()

	mockDB.
		On("PingContext", mock.Anything).
		Return(nil).
		Once()

	svc := &PosService{
		db:        mockDB,
		validator: mockValidator,
	}

	err := svc.GetHealthByTenantID(context.Background(), "tenant_001")

	assert.NoError(t, err)

	mockValidator.AssertExpectations(t)
	mockDB.AssertExpectations(t)
}

func TestGetBranchesByTenantID_Success(t *testing.T) {
	repo := &MockRepo{
		data: []domain.BranchResponse{
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
		},
	}
	svc := NewPosService(nil, repo)

	got, err := svc.GetBranchesByTenantID(context.Background(), "aura-bkk")
	assert.NoError(t, err)
	assert.Len(t, got, 2)
	assert.Equal(t, "bkk-001", got[0].BranchID)
}

func TestGetBranchesByTenantID_InvalidStatus(t *testing.T) {
	repo := &MockRepo{
		data: []domain.BranchResponse{
			{
				BranchID:   "bkk-001",
				BranchName: "Aura Siam",
				Status:     "pending",
			},
		},
	}

	svc := NewPosService(nil, repo)

	got, err := svc.GetBranchesByTenantID(context.Background(), "aura-bkk")

	assert.Error(t, err)
	assert.Nil(t, got)
	assert.Equal(t, "status must be active or inactive", err.Error())
}