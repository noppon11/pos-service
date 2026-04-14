package repository

import (
	"context"
	"database/sql"
	"os"
	"testing"
	appErr "pos-service/internal/errors"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/joho/godotenv"
)

func setupTestPostgresDB(t *testing.T) *sql.DB {
	t.Helper()

	_ = godotenv.Load(".env.test")

	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}

	db, err := sql.Open("postgres", dbURL)
	require.NoError(t, err)

	err = db.Ping()
	require.NoError(t, err)

	resetQuery := `
	DROP TABLE IF EXISTS branches;

	CREATE TABLE branches (
		id BIGSERIAL PRIMARY KEY,
		tenant_id VARCHAR(100) NOT NULL,
		branch_id VARCHAR(100) NOT NULL,
		branch_name VARCHAR(255) NOT NULL,
		status VARCHAR(50) NOT NULL,
		timezone VARCHAR(100) NOT NULL,
		currency VARCHAR(20) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
		UNIQUE (tenant_id, branch_id)
	);
	`
	_, err = db.Exec(resetQuery)
	require.NoError(t, err)

	seedQuery := `
	INSERT INTO branches (tenant_id, branch_id, branch_name, status, timezone, currency)
	VALUES
	('aura-bkk', 'bkk-001', 'Aura Siam', 'active', 'Asia/Bangkok', 'THB'),
	('aura-bkk', 'bkk-002', 'Aura Ari', 'inactive', 'Asia/Bangkok', 'THB'),
	('aura-cnx', 'cnx-001', 'Aura Chiang Mai', 'active', 'Asia/Bangkok', 'THB');
	`
	_, err = db.Exec(seedQuery)
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = db.Close()
	})

	return db
}

func TestPostgresBranchRepository_ListByTenantID_Success(t *testing.T) {
	db := setupTestPostgresDB(t)
	repo := NewPostgresBranchRepository(db)

	got, err := repo.ListByTenantID(context.Background(), "aura-bkk")

	require.NoError(t, err)
	require.Len(t, got, 2)

	assert.Equal(t, "bkk-002", got[0].BranchID)
	assert.Equal(t, "Aura Ari", got[0].BranchName)
	assert.Equal(t, "inactive", got[0].Status)
	assert.Equal(t, "Asia/Bangkok", got[0].Timezone)
	assert.Equal(t, "THB", got[0].Currency)

	assert.Equal(t, "bkk-001", got[1].BranchID)
	assert.Equal(t, "Aura Siam", got[1].BranchName)
	assert.Equal(t, "active", got[1].Status)
	assert.Equal(t, "Asia/Bangkok", got[1].Timezone)
	assert.Equal(t, "THB", got[1].Currency)
}

func TestPostgresBranchRepository_ListByTenantID_EmptyResult(t *testing.T) {
	db := setupTestPostgresDB(t)
	repo := NewPostgresBranchRepository(db)

	got, err := repo.ListByTenantID(context.Background(), "unknown-tenant")

	require.NoError(t, err)
	assert.Len(t, got, 0)
}

func TestPostgresBranchRepository_GetByTenantIDAndBranchID_Success(t *testing.T) {
	db := setupTestPostgresDB(t)
	repo := NewPostgresBranchRepository(db)

	got, err := repo.GetByTenantIDAndBranchID(context.Background(), "aura-bkk", "bkk-001")

	require.NoError(t, err)
	require.NotNil(t, got)

	assert.Equal(t, "bkk-001", got.BranchID)
	assert.Equal(t, "Aura Siam", got.BranchName)
	assert.Equal(t, "active", got.Status)
	assert.Equal(t, "Asia/Bangkok", got.Timezone)
	assert.Equal(t, "THB", got.Currency)
}

func TestPostgresBranchRepository_GetByTenantIDAndBranchID_TenantNotFound(t *testing.T) {
	db := setupTestPostgresDB(t)
	repo := NewPostgresBranchRepository(db)

	got, err := repo.GetByTenantIDAndBranchID(context.Background(), "unknown-tenant", "bkk-001")

	assert.Nil(t, got)
	assert.ErrorIs(t, err, appErr.ErrTenantNotFound)
}

func TestPostgresBranchRepository_GetByTenantIDAndBranchID_BranchNotFound(t *testing.T) {
	db := setupTestPostgresDB(t)
	repo := NewPostgresBranchRepository(db)

	got, err := repo.GetByTenantIDAndBranchID(context.Background(), "aura-bkk", "not-found")

	assert.Nil(t, got)
	assert.ErrorIs(t, err, appErr.ErrBranchNotFound)
}