package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInMemoryBranchRepository_ListByTenantID(t *testing.T) {
	repo := NewInMemoryBranchRepository()

	t.Run("found tenant", func(t *testing.T) {
		data, err := repo.ListByTenantID(context.Background(), "aura-bkk")
		assert.NoError(t, err)
		assert.Len(t, data, 2)
	})

	t.Run("not found tenant", func(t *testing.T) {
		data, err := repo.ListByTenantID(context.Background(), "unknown")
		assert.NoError(t, err)
		assert.Len(t, data, 0)
	})
}