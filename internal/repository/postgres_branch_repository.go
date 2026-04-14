package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"pos-service/internal/domain"
	appErr "pos-service/internal/errors"
)

type PostgresBranchRepository struct {
	db *sql.DB
}

func NewPostgresBranchRepository(db *sql.DB) *PostgresBranchRepository {
	return &PostgresBranchRepository{
		db: db,
	}
}

func (r *PostgresBranchRepository) ListByTenantID(ctx context.Context, tenantID string) ([]domain.BranchResponse, error) {
	const query = `
		SELECT
			branch_id,
			branch_name,
			status,
			timezone,
			currency
		FROM branches
		WHERE tenant_id = $1
		ORDER BY branch_name ASC
	`

	rows, err := r.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, fmt.Errorf("query branches by tenant id: %w", err)
	}
	defer rows.Close()

	branches := make([]domain.BranchResponse, 0)
	for rows.Next() {
		var branch domain.BranchResponse

		if err := rows.Scan(
			&branch.BranchID,
			&branch.BranchName,
			&branch.Status,
			&branch.Timezone,
			&branch.Currency,
		); err != nil {
			return nil, fmt.Errorf("scan branch row: %w", err)
		}

		branches = append(branches, branch)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate branch rows: %w", err)
	}

	// keep same behavior as your mock repo:
	// if tenant has no branches, return empty list, nil
	return branches, nil
}

func (r *PostgresBranchRepository) GetByTenantIDAndBranchID(ctx context.Context, tenantID, branchID string) (*domain.BranchResponse, error) {
	const tenantCheckQuery = `
		SELECT EXISTS (
			SELECT 1
			FROM branches
			WHERE tenant_id = $1
		)
	`

	var tenantExists bool
	if err := r.db.QueryRowContext(ctx, tenantCheckQuery, tenantID).Scan(&tenantExists); err != nil {
		return nil, fmt.Errorf("check tenant existence: %w", err)
	}

	if !tenantExists {
		return nil, appErr.ErrTenantNotFound
	}

	const query = `
		SELECT
			branch_id,
			branch_name,
			status,
			timezone,
			currency
		FROM branches
		WHERE tenant_id = $1
		  AND branch_id = $2
		LIMIT 1
	`

	var branch domain.BranchResponse
	err := r.db.QueryRowContext(ctx, query, tenantID, branchID).Scan(
		&branch.BranchID,
		&branch.BranchName,
		&branch.Status,
		&branch.Timezone,
		&branch.Currency,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, appErr.ErrBranchNotFound
		}
		return nil, fmt.Errorf("get branch by tenant id and branch id: %w", err)
	}

	return &branch, nil
}