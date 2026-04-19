package repository

import (
	"context"
	"database/sql"
	"errors"

	"pos-service/internal/domain"
)

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	const q = `
		SELECT id, email, password_hash, full_name, tenant_id, role, is_active, created_at, updated_at
		FROM users
		WHERE email = $1
		LIMIT 1
	`

	var u domain.User
	err := r.db.QueryRowContext(ctx, q, email).Scan(
		&u.ID,
		&u.Email,
		&u.PasswordHash,
		&u.FullName,
		&u.TenantID,
		&u.Role,
		&u.IsActive,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	branchIDs, err := r.getUserBranchIDs(ctx, u.ID)
	if err != nil {
		return nil, err
	}
	u.BranchIDs = branchIDs

	return &u, nil
}

func (r *PostgresUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	const q = `
		SELECT id, email, password_hash, full_name, tenant_id, role, is_active, created_at, updated_at
		FROM users
		WHERE id = $1
		LIMIT 1
	`

	var u domain.User
	err := r.db.QueryRowContext(ctx, q, id).Scan(
		&u.ID,
		&u.Email,
		&u.PasswordHash,
		&u.FullName,
		&u.TenantID,
		&u.Role,
		&u.IsActive,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	branchIDs, err := r.getUserBranchIDs(ctx, u.ID)
	if err != nil {
		return nil, err
	}
	u.BranchIDs = branchIDs

	return &u, nil
}

func (r *PostgresUserRepository) getUserBranchIDs(ctx context.Context, userID string) ([]string, error) {
	const q = `
		SELECT branch_id
		FROM user_branch_access
		WHERE user_id = $1
		ORDER BY branch_id
	`

	rows, err := r.db.QueryContext(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var branchIDs []string
	for rows.Next() {
		var branchID string
		if err := rows.Scan(&branchID); err != nil {
			return nil, err
		}
		branchIDs = append(branchIDs, branchID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return branchIDs, nil
}