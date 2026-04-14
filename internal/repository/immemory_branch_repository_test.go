package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	appErr "pos-service/internal/errors"
)

func TestNewInMemoryBranchRepository(t *testing.T) {
	repo := NewInMemoryBranchRepository()

	assert.NotNil(t, repo)
	assert.NotNil(t, repo.data)
	assert.Contains(t, repo.data, "aura-bkk")
	assert.Contains(t, repo.data, "aura-cnx")
}

func TestInMemoryBranchRepository_ListByTenantID_Success(t *testing.T) {
	repo := NewInMemoryBranchRepository()

	got, err := repo.ListByTenantID(context.Background(), "aura-bkk")

	assert.NoError(t, err)
	assert.Len(t, got, 2)
	assert.Equal(t, "bkk-001", got[0].BranchID)
	assert.Equal(t, "Aura Siam", got[0].BranchName)
	assert.Equal(t, "active", got[0].Status)
	assert.Equal(t, "Asia/Bangkok", got[0].Timezone)
	assert.Equal(t, "THB", got[0].Currency)
}

func TestInMemoryBranchRepository_ListByTenantID_TenantNotFound(t *testing.T) {
	repo := NewInMemoryBranchRepository()

	got, err := repo.ListByTenantID(context.Background(), "unknown-tenant")

	assert.NoError(t, err)
	assert.NotNil(t, got)
	assert.Len(t, got, 0)
}

func TestInMemoryBranchRepository_GetByTenantIDAndBranchID_Success(t *testing.T) {
	repo := NewInMemoryBranchRepository()

	got, err := repo.GetByTenantIDAndBranchID(context.Background(), "aura-bkk", "bkk-001")

	assert.NoError(t, err)
	assert.NotNil(t, got)
	assert.Equal(t, "bkk-001", got.BranchID)
	assert.Equal(t, "Aura Siam", got.BranchName)
	assert.Equal(t, "active", got.Status)
	assert.Equal(t, "Asia/Bangkok", got.Timezone)
	assert.Equal(t, "THB", got.Currency)
}

func TestInMemoryBranchRepository_GetByTenantIDAndBranchID_TenantNotFound(t *testing.T) {
	repo := NewInMemoryBranchRepository()

	got, err := repo.GetByTenantIDAndBranchID(context.Background(), "unknown-tenant", "bkk-001")

	assert.Nil(t, got)
	assert.ErrorIs(t, err, appErr.ErrTenantNotFound)
}

func TestInMemoryBranchRepository_GetByTenantIDAndBranchID_BranchNotFound(t *testing.T) {
	repo := NewInMemoryBranchRepository()

	got, err := repo.GetByTenantIDAndBranchID(context.Background(), "aura-bkk", "not-found")

	assert.Nil(t, got)
	assert.ErrorIs(t, err, appErr.ErrBranchNotFound)
}