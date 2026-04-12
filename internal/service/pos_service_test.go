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

func (m *MockDB) PingContext(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
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

type MockValidator struct {
	mock.Mock
}

func (m *MockValidator) TenantIDValidation(tenantID string) error {
	args := m.Called(tenantID)
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
	mockValidator := new(MockValidator)

	mockValidator.
		On("TenantIDValidation", "aura-bkk").
		Return(nil).
		Once()

	svc := &PosService{
		branchRepo: repo,
		validator:  mockValidator,
	}

	got, err := svc.GetBranchesByTenantID(context.Background(), "aura-bkk")

	assert.NoError(t, err)
	assert.Len(t, got, 2)
	assert.Equal(t, "bkk-001", got[0].BranchID)
	assert.Equal(t, "Aura Siam", got[0].BranchName)
	assert.Equal(t, "active", got[0].Status)
	mockValidator.AssertExpectations(t)
}

func TestGetBranchesByTenantID_InvalidTenantID(t *testing.T) {
	repo := &MockRepo{}
	mockValidator := new(MockValidator)

	mockValidator.
		On("TenantIDValidation", "AURA").
		Return(errors.New("tenant_id must be 3-50 chars, lowercase letters, numbers, underscore or dash only")).
		Once()

	svc := &PosService{
		branchRepo: repo,
		validator:  mockValidator,
	}

	got, err := svc.GetBranchesByTenantID(context.Background(), "AURA")

	assert.Error(t, err)
	assert.Nil(t, got)
	assert.EqualError(t, err, "tenant_id must be 3-50 chars, lowercase letters, numbers, underscore or dash only")
	mockValidator.AssertExpectations(t)
}

func TestGetBranchesByTenantID_RepoError(t *testing.T) {
	repo := &MockRepo{
		err: errors.New("repository error"),
	}
	mockValidator := new(MockValidator)

	mockValidator.
		On("TenantIDValidation", "aura-bkk").
		Return(nil).
		Once()

	svc := &PosService{
		branchRepo: repo,
		validator:  mockValidator,
	}

	got, err := svc.GetBranchesByTenantID(context.Background(), "aura-bkk")

	assert.Error(t, err)
	assert.Nil(t, got)
	assert.EqualError(t, err, "repository error")
	mockValidator.AssertExpectations(t)
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
	mockValidator := new(MockValidator)

	mockValidator.
		On("TenantIDValidation", "aura-bkk").
		Return(nil).
		Once()

	svc := &PosService{
		branchRepo: repo,
		validator:  mockValidator,
	}

	got, err := svc.GetBranchesByTenantID(context.Background(), "aura-bkk")

	assert.Error(t, err)
	assert.Nil(t, got)
	assert.EqualError(t, err, "status must be active or inactive")
	mockValidator.AssertExpectations(t)
}

func TestGetBranchesByTenantID_EmptyBranchID(t *testing.T) {
	repo := &MockRepo{
		data: []domain.BranchResponse{
			{
				BranchID:   "",
				BranchName: "Aura Siam",
				Status:     "active",
			},
		},
	}
	mockValidator := new(MockValidator)

	mockValidator.
		On("TenantIDValidation", "aura-bkk").
		Return(nil).
		Once()

	svc := &PosService{
		branchRepo: repo,
		validator:  mockValidator,
	}

	got, err := svc.GetBranchesByTenantID(context.Background(), "aura-bkk")

	assert.Error(t, err)
	assert.Nil(t, got)
	assert.EqualError(t, err, "branch_id is required")
	mockValidator.AssertExpectations(t)
}

func TestGetBranchesByTenantID_EmptyBranchName(t *testing.T) {
	repo := &MockRepo{
		data: []domain.BranchResponse{
			{
				BranchID:   "bkk-001",
				BranchName: "",
				Status:     "active",
			},
		},
	}
	mockValidator := new(MockValidator)

	mockValidator.
		On("TenantIDValidation", "aura-bkk").
		Return(nil).
		Once()

	svc := &PosService{
		branchRepo: repo,
		validator:  mockValidator,
	}

	got, err := svc.GetBranchesByTenantID(context.Background(), "aura-bkk")

	assert.Error(t, err)
	assert.Nil(t, got)
	assert.EqualError(t, err, "branch_name is required")
	mockValidator.AssertExpectations(t)
}