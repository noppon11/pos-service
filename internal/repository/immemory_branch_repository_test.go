package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInMemoryBranchRepository_ListByTenantID(t *testing.T) {
	repo := NewInMemoryBranchRepository()

	t.Run("found tenant aura-bkk", func(t *testing.T) {
		data, err := repo.ListByTenantID(context.Background(), "aura-bkk")

		assert.NoError(t, err)
		assert.Len(t, data, 2)

		assert.Equal(t, "bkk-001", data[0].BranchID)
		assert.Equal(t, "Aura Siam", data[0].BranchName)
		assert.Equal(t, "active", data[0].Status)

		assert.Equal(t, "bkk-002", data[1].BranchID)
		assert.Equal(t, "Aura Ari", data[1].BranchName)
		assert.Equal(t, "inactive", data[1].Status)
	})

	t.Run("found tenant aura-cnx", func(t *testing.T) {
		data, err := repo.ListByTenantID(context.Background(), "aura-cnx")

		assert.NoError(t, err)
		assert.Len(t, data, 1)

		assert.Equal(t, "cnx-001", data[0].BranchID)
		assert.Equal(t, "Aura Chiang Mai", data[0].BranchName)
		assert.Equal(t, "active", data[0].Status)
	})

	t.Run("unknown tenant returns empty list", func(t *testing.T) {
		data, err := repo.ListByTenantID(context.Background(), "unknown")

		assert.NoError(t, err)
		assert.NotNil(t, data)
		assert.Len(t, data, 0)
	})
}

func TestInMemoryBranchRepository_GetByTenantIDAndBranchID(t *testing.T) {
	repo := NewInMemoryBranchRepository()

	t.Run("found branch in tenant", func(t *testing.T) {
		data, err := repo.GetByTenantIDAndBranchID(context.Background(), "aura-bkk", "bkk-001")

		assert.NoError(t, err)
		assert.NotNil(t, data)
		assert.Equal(t, "bkk-001", data.BranchID)
		assert.Equal(t, "Aura Siam", data.BranchName)
		assert.Equal(t, "active", data.Status)
	})

	t.Run("found another branch in tenant", func(t *testing.T) {
		data, err := repo.GetByTenantIDAndBranchID(context.Background(), "aura-bkk", "bkk-002")

		assert.NoError(t, err)
		assert.NotNil(t, data)
		assert.Equal(t, "bkk-002", data.BranchID)
		assert.Equal(t, "Aura Ari", data.BranchName)
		assert.Equal(t, "inactive", data.Status)
	})

	t.Run("unknown tenant returns ErrTenantNotFound", func(t *testing.T) {
		data, err := repo.GetByTenantIDAndBranchID(context.Background(), "unknown", "bkk-001")

		assert.Nil(t, data)
		assert.ErrorIs(t, err, ErrTenantNotFound)
	})

	t.Run("branch not found in existing tenant returns ErrBranchNotFound", func(t *testing.T) {
		data, err := repo.GetByTenantIDAndBranchID(context.Background(), "aura-bkk", "bkk-999")

		assert.Nil(t, data)
		assert.ErrorIs(t, err, ErrBranchNotFound)
	})
}