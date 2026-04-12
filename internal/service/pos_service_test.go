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

func TestGetHealth_DBNotConfigured(t *testing.T) {
	svc := NewPosService(nil, nil)

	err := svc.GetHealth(context.Background())

	assert.ErrorIs(t, err, ErrDBNotConfigured)
}

func TestGetHealth_DBDown(t *testing.T) {
	mockDB := new(MockDB)

	mockDB.
		On("PingContext", mock.Anything).
		Return(errors.New("db down")).
		Once()

	svc := NewPosService(mockDB, nil)

	err := svc.GetHealth(context.Background())

	assert.EqualError(t, err, "db down")
	mockDB.AssertExpectations(t)
}

func TestGetHealth_Success(t *testing.T) {
	mockDB := new(MockDB)

	mockDB.
		On("PingContext", mock.Anything).
		Return(nil).
		Once()

	svc := NewPosService(mockDB, nil)

	err := svc.GetHealth(context.Background())

	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
}

func TestGetHealthByTenantID_DBNotConfigured(t *testing.T) {
	svc := NewPosService(nil, nil)

	err := svc.GetHealthByTenantID(context.Background(), "tenant_001")

	assert.ErrorIs(t, err, ErrDBNotConfigured)
}

func TestGetHealthByTenantID_DBDown(t *testing.T) {
	mockDB := new(MockDB)

	mockDB.
		On("PingContext", mock.Anything).
		Return(errors.New("db down")).
		Once()

	svc := NewPosService(mockDB, nil)

	err := svc.GetHealthByTenantID(context.Background(), "tenant_001")

	assert.EqualError(t, err, "db down")
	mockDB.AssertExpectations(t)
}

func TestGetHealthByTenantID_Success(t *testing.T) {
	mockDB := new(MockDB)

	mockDB.
		On("PingContext", mock.Anything).
		Return(nil).
		Once()

	svc := NewPosService(mockDB, nil)

	err := svc.GetHealthByTenantID(context.Background(), "tenant_001")

	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
}

func TestGetBranchesByTenantID_RepoNotConfigured(t *testing.T) {
	svc := NewPosService(nil, nil)

	got, err := svc.GetBranchesByTenantID(context.Background(), "aura-bkk")

	assert.Nil(t, got)
	assert.ErrorIs(t, err, ErrBranchRepoNotConfigured)
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
	assert.Equal(t, "Aura Siam", got[0].BranchName)
	assert.Equal(t, "active", got[0].Status)
}

func TestGetBranchesByTenantID_RepoError(t *testing.T) {
	repo := &MockRepo{
		err: errors.New("repository error"),
	}

	svc := NewPosService(nil, repo)

	got, err := svc.GetBranchesByTenantID(context.Background(), "aura-bkk")

	assert.Nil(t, got)
	assert.EqualError(t, err, "repository error")
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

	assert.Nil(t, got)
	assert.ErrorIs(t, err, ErrInvalidBranchStatus)
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

	svc := NewPosService(nil, repo)

	got, err := svc.GetBranchesByTenantID(context.Background(), "aura-bkk")

	assert.Nil(t, got)
	assert.ErrorIs(t, err, ErrBranchIDRequired)
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

	svc := NewPosService(nil, repo)

	got, err := svc.GetBranchesByTenantID(context.Background(), "aura-bkk")

	assert.Nil(t, got)
	assert.ErrorIs(t, err, ErrBranchNameRequired)
}