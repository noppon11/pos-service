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
	data map[string][]domain.BranchResponse
	err  error
}

func (m *MockRepo) ListByTenantID(ctx context.Context, tenantID string) ([]domain.BranchResponse, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.data[tenantID], nil
}

func (m *MockRepo) GetByTenantIDAndBranchID(ctx context.Context, tenantID string, branchID string) (*domain.BranchResponse, error) {
	if m.err != nil {
		return nil, m.err
	}

	branches, ok := m.data[tenantID]
	if !ok {
		return nil, nil
	}

	for i := range branches {
		if branches[i].BranchID == branchID {
			return &branches[i], nil
		}
	}

	return nil, nil
}

type MockValidator struct {
	mock.Mock
}

func (m *MockValidator) BranchValidation(branch domain.BranchResponse) error {
	args := m.Called(branch)
	return args.Error(0)
}

func TestGetHealth_DBNotConfigured(t *testing.T) {
	svc := NewPosService(nil, nil, nil)

	err := svc.GetHealth(context.Background())

	assert.ErrorIs(t, err, ErrDBNotConfigured)
}

func TestGetHealth_DBDown(t *testing.T) {
	mockDB := new(MockDB)

	mockDB.On("PingContext", mock.Anything).Return(errors.New("db down")).Once()

	svc := NewPosService(mockDB, nil, nil)

	err := svc.GetHealth(context.Background())

	assert.EqualError(t, err, "db down")
	mockDB.AssertExpectations(t)
}

func TestGetHealth_Success(t *testing.T) {
	mockDB := new(MockDB)

	mockDB.On("PingContext", mock.Anything).Return(nil).Once()

	svc := NewPosService(mockDB, nil, nil)

	err := svc.GetHealth(context.Background())

	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
}

func TestGetHealthByTenantID_DBNotConfigured(t *testing.T) {
	svc := NewPosService(nil, nil, nil)

	err := svc.GetHealthByTenantID(context.Background(), "tenant_001")

	assert.ErrorIs(t, err, ErrDBNotConfigured)
}

func TestGetHealthByTenantID_DBDown(t *testing.T) {
	mockDB := new(MockDB)

	mockDB.On("PingContext", mock.Anything).Return(errors.New("db down")).Once()

	svc := NewPosService(mockDB, nil, nil)

	err := svc.GetHealthByTenantID(context.Background(), "tenant_001")

	assert.EqualError(t, err, "db down")
	mockDB.AssertExpectations(t)
}

func TestGetHealthByTenantID_Success(t *testing.T) {
	mockDB := new(MockDB)

	mockDB.On("PingContext", mock.Anything).Return(nil).Once()

	svc := NewPosService(mockDB, nil, nil)

	err := svc.GetHealthByTenantID(context.Background(), "tenant_001")

	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
}

func TestGetBranchesByTenantID_RepoNotConfigured(t *testing.T) {
	svc := NewPosService(nil, nil, nil)

	got, err := svc.GetBranchesByTenantID(context.Background(), "aura-bkk")

	assert.Nil(t, got)
	assert.ErrorIs(t, err, ErrBranchRepoNotConfigured)
}

func TestGetBranchesByTenantID_ValidatorNotConfigured(t *testing.T) {
	repo := &MockRepo{
		data: map[string][]domain.BranchResponse{
			"aura-bkk": {
				{
					BranchID:   "bkk-001",
					BranchName: "Aura Siam",
					Status:     "active",
					Timezone:   "Asia/Bangkok",
					Currency:   "THB",
				},
			},
		},
	}

	svc := NewPosService(nil, repo, nil)

	got, err := svc.GetBranchesByTenantID(context.Background(), "aura-bkk")

	assert.Nil(t, got)
	assert.ErrorIs(t, err, ErrValidatorNotConfigured)
}

func TestGetBranchesByTenantID_Success(t *testing.T) {
	repo := &MockRepo{
		data: map[string][]domain.BranchResponse{
			"aura-bkk": {
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
			},
		},
	}

	validator := new(MockValidator)
	validator.On("BranchValidation", mock.Anything).Return(nil).Twice()

	svc := NewPosService(nil, repo, validator)

	got, err := svc.GetBranchesByTenantID(context.Background(), "aura-bkk")

	assert.NoError(t, err)
	assert.Len(t, got, 2)
	assert.Equal(t, "bkk-001", got[0].BranchID)
	assert.Equal(t, "Aura Siam", got[0].BranchName)
	assert.Equal(t, "active", got[0].Status)
	validator.AssertExpectations(t)
}

func TestGetBranchesByTenantID_RepoError(t *testing.T) {
	repo := &MockRepo{
		err: errors.New("repository error"),
	}
	validator := new(MockValidator)

	svc := NewPosService(nil, repo, validator)

	got, err := svc.GetBranchesByTenantID(context.Background(), "aura-bkk")

	assert.Nil(t, got)
	assert.EqualError(t, err, "repository error")
}

func TestGetBranchesByTenantID_InvalidStatus(t *testing.T) {
	repo := &MockRepo{
		data: map[string][]domain.BranchResponse{
			"aura-bkk": {
				{
					BranchID:   "bkk-001",
					BranchName: "Aura Siam",
					Status:     "pending",
					Timezone:   "Asia/Bangkok",
					Currency:   "THB",
				},
			},
		},
	}

	validator := new(MockValidator)
	validator.On("BranchValidation", mock.Anything).Return(ErrInvalidBranchStatus).Once()

	svc := NewPosService(nil, repo, validator)

	got, err := svc.GetBranchesByTenantID(context.Background(), "aura-bkk")

	assert.Nil(t, got)
	assert.ErrorIs(t, err, ErrInvalidBranchStatus)
	validator.AssertExpectations(t)
}

func TestGetBranchesByTenantID_EmptyBranchID(t *testing.T) {
	repo := &MockRepo{
		data: map[string][]domain.BranchResponse{
			"aura-bkk": {
				{
					BranchID:   "",
					BranchName: "Aura Siam",
					Status:     "active",
					Timezone:   "Asia/Bangkok",
					Currency:   "THB",
				},
			},
		},
	}

	validator := new(MockValidator)
	validator.On("BranchValidation", mock.Anything).Return(ErrBranchIDRequired).Once()

	svc := NewPosService(nil, repo, validator)

	got, err := svc.GetBranchesByTenantID(context.Background(), "aura-bkk")

	assert.Nil(t, got)
	assert.ErrorIs(t, err, ErrBranchIDRequired)
	validator.AssertExpectations(t)
}

func TestGetBranchesByTenantID_EmptyBranchName(t *testing.T) {
	repo := &MockRepo{
		data: map[string][]domain.BranchResponse{
			"aura-bkk": {
				{
					BranchID:   "bkk-001",
					BranchName: "",
					Status:     "active",
					Timezone:   "Asia/Bangkok",
					Currency:   "THB",
				},
			},
		},
	}

	validator := new(MockValidator)
	validator.On("BranchValidation", mock.Anything).Return(ErrBranchNameRequired).Once()

	svc := NewPosService(nil, repo, validator)

	got, err := svc.GetBranchesByTenantID(context.Background(), "aura-bkk")

	assert.Nil(t, got)
	assert.ErrorIs(t, err, ErrBranchNameRequired)
	validator.AssertExpectations(t)
}

func TestGetBranchDetail_RepoNotConfigured(t *testing.T) {
	svc := NewPosService(nil, nil, nil)

	got, err := svc.GetBranchDetail(context.Background(), "aura-bkk", "bkk-001")

	assert.Nil(t, got)
	assert.ErrorIs(t, err, ErrBranchRepoNotConfigured)
}

func TestGetBranchDetail_ValidatorNotConfigured(t *testing.T) {
	repo := &MockRepo{
		data: map[string][]domain.BranchResponse{
			"aura-bkk": {
				{
					BranchID:   "bkk-001",
					BranchName: "Aura Siam",
					Status:     "active",
					Timezone:   "Asia/Bangkok",
					Currency:   "THB",
				},
			},
		},
	}

	svc := NewPosService(nil, repo, nil)

	got, err := svc.GetBranchDetail(context.Background(), "aura-bkk", "bkk-001")

	assert.Nil(t, got)
	assert.ErrorIs(t, err, ErrValidatorNotConfigured)
}

func TestGetBranchDetail_RepoError(t *testing.T) {
	repo := &MockRepo{
		err: errors.New("repository error"),
	}
	validator := new(MockValidator)

	svc := NewPosService(nil, repo, validator)

	got, err := svc.GetBranchDetail(context.Background(), "aura-bkk", "bkk-001")

	assert.Nil(t, got)
	assert.EqualError(t, err, "repository error")
}

func TestGetBranchDetail_Success(t *testing.T) {
	repo := &MockRepo{
		data: map[string][]domain.BranchResponse{
			"aura-bkk": {
				{
					BranchID:   "bkk-001",
					BranchName: "Aura Siam",
					Status:     "active",
					Timezone:   "Asia/Bangkok",
					Currency:   "THB",
				},
			},
		},
	}

	validator := new(MockValidator)
	validator.On("BranchValidation", mock.Anything).Return(nil).Once()

	svc := NewPosService(nil, repo, validator)

	got, err := svc.GetBranchDetail(context.Background(), "aura-bkk", "bkk-001")

	assert.NoError(t, err)
	assert.NotNil(t, got)
	assert.Equal(t, "bkk-001", got.BranchID)
	assert.Equal(t, "Aura Siam", got.BranchName)
	assert.Equal(t, "active", got.Status)
	assert.Equal(t, "Asia/Bangkok", got.Timezone)
	assert.Equal(t, "THB", got.Currency)
	validator.AssertExpectations(t)
}

func TestGetBranchDetail_InvalidStatus(t *testing.T) {
	repo := &MockRepo{
		data: map[string][]domain.BranchResponse{
			"aura-bkk": {
				{
					BranchID:   "bkk-001",
					BranchName: "Aura Siam",
					Status:     "pending",
					Timezone:   "Asia/Bangkok",
					Currency:   "THB",
				},
			},
		},
	}

	validator := new(MockValidator)
	validator.On("BranchValidation", mock.Anything).Return(ErrInvalidBranchStatus).Once()

	svc := NewPosService(nil, repo, validator)

	got, err := svc.GetBranchDetail(context.Background(), "aura-bkk", "bkk-001")

	assert.Nil(t, got)
	assert.ErrorIs(t, err, ErrInvalidBranchStatus)
	validator.AssertExpectations(t)
}

func TestGetBranchDetail_EmptyBranchID(t *testing.T) {
	repo := &MockRepo{
		data: map[string][]domain.BranchResponse{
			"aura-bkk": {
				{
					BranchID:   "",
					BranchName: "Aura_Siam",
					Status:     "active",
					Timezone:   "Asia/Bangkok",
					Currency:   "THB",
				},
			},
		},
	}

	validator := new(MockValidator)
	validator.On("BranchValidation", mock.Anything).Return(ErrBranchIDRequired).Once()

	svc := NewPosService(nil, repo, validator)

	got, err := svc.GetBranchDetail(context.Background(), "aura-bkk", "")

	assert.Error(t, err)
	assert.Nil(t, got)
}

func TestGetBranchDetail_EmptyBranchName(t *testing.T) {
	repo := &MockRepo{
		data: map[string][]domain.BranchResponse{
			"aura-bkk": {
				{
					BranchID:   "bkk-001",
					BranchName: "",
					Status:     "active",
					Timezone:   "Asia/Bangkok",
					Currency:   "THB",
				},
			},
		},
	}

	validator := new(MockValidator)
	validator.On("BranchValidation", mock.Anything).Return(ErrBranchNameRequired).Once()

	svc := NewPosService(nil, repo, validator)

	got, err := svc.GetBranchDetail(context.Background(), "aura-bkk", "bkk-001")

	assert.Nil(t, got)
	assert.ErrorIs(t, err, ErrBranchNameRequired)
	validator.AssertExpectations(t)
}