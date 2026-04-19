package repository

import (
	"context"
	"database/sql"
	"errors"

	"pos-service/internal/domain"
	appErr "pos-service/internal/errors"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type PostgresProductRepository struct {
	db *sql.DB
}

func NewPostgresProductRepository(db *sql.DB) *PostgresProductRepository {
	return &PostgresProductRepository{db: db}
}

func (r *PostgresProductRepository) ListByTenantIDAndBranchID(
	ctx context.Context,
	tenantID string,
	branchID string,
	filter ProductListFilter,
) ([]domain.Product, int, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 20
	}

	offset := (filter.Page - 1) * filter.Limit

	countQuery := `
		SELECT COUNT(*)
		FROM products
		WHERE tenant_id = $1
		  AND branch_id = $2
		  AND deleted_at IS NULL
		  AND ($3 = '' OR category_id = $3)
	`

	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, tenantID, branchID, filter.CategoryID).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT product_id, name, sku, price, category_id, unit, is_active, deleted_at
		FROM products
		WHERE tenant_id = $1
		  AND branch_id = $2
		  AND deleted_at IS NULL
		  AND ($3 = '' OR category_id = $3)
		ORDER BY created_at DESC
		LIMIT $4 OFFSET $5
	`

	rows, err := r.db.QueryContext(ctx, query, tenantID, branchID, filter.CategoryID, filter.Limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		var p domain.Product
		var deletedAt sql.NullTime

		if err := rows.Scan(
			&p.ProductID,
			&p.Name,
			&p.SKU,
			&p.Price,
			&p.CategoryID,
			&p.Unit,
			&p.IsActive,
			&deletedAt,
		); err != nil {
			return nil, 0, err
		}

		if deletedAt.Valid {
			p.DeletedAt = &deletedAt.Time
		}

		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *PostgresProductRepository) GetByTenantIDBranchIDAndProductID(
	ctx context.Context,
	tenantID string,
	branchID string,
	productID string,
) (*domain.Product, error) {
	query := `
		SELECT product_id, name, sku, price, category_id, unit, is_active, deleted_at
		FROM products
		WHERE tenant_id = $1 AND branch_id = $2 AND product_id = $3 AND deleted_at IS NULL
	`

	var p domain.Product
	var deletedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, tenantID, branchID, productID).Scan(
		&p.ProductID,
		&p.Name,
		&p.SKU,
		&p.Price,
		&p.CategoryID,
		&p.Unit,
		&p.IsActive,
		&deletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if deletedAt.Valid {
		p.DeletedAt = &deletedAt.Time
	}

	return &p, nil
}

func (r *PostgresProductRepository) Create(
	ctx context.Context,
	tenantID string,
	branchID string,
	product domain.Product,
) (*domain.Product, error) {
	product.ProductID = uuid.New().String()

	query := `
		INSERT INTO products (
			product_id,
			tenant_id,
			branch_id,
			name,
			sku,
			price,
			category_id,
			unit,
			is_active
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING product_id, name, sku, price, category_id, unit, is_active
	`

	var created domain.Product
	err := r.db.QueryRowContext(
		ctx,
		query,
		product.ProductID,
		tenantID,
		branchID,
		product.Name,
		product.SKU,
		product.Price,
		product.CategoryID,
		product.Unit,
		product.IsActive,
	).Scan(
		&created.ProductID,
		&created.Name,
		&created.SKU,
		&created.Price,
		&created.CategoryID,
		&created.Unit,
		&created.IsActive,
	)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return nil, appErr.ErrProductAlreadyExists
		}
		return nil, err
	}

	return &created, nil
}

func (r *PostgresProductRepository) Update(
	ctx context.Context,
	tenantID string,
	branchID string,
	productID string,
	product domain.Product,
) (*domain.Product, error) {
	query := `
		UPDATE products
		SET name = $4,
		    sku = $5,
		    price = $6,
		    category_id = $7,
		    unit = $8,
		    is_active = $9,
		    updated_at = NOW()
		WHERE tenant_id = $1
		  AND branch_id = $2
		  AND product_id = $3
		  AND deleted_at IS NULL
		RETURNING product_id, name, sku, price, category_id, unit, is_active, deleted_at
	`

	var updated domain.Product
	var deletedAt sql.NullTime

	err := r.db.QueryRowContext(
		ctx,
		query,
		tenantID,
		branchID,
		productID,
		product.Name,
		product.SKU,
		product.Price,
		product.CategoryID,
		product.Unit,
		product.IsActive,
	).Scan(
		&updated.ProductID,
		&updated.Name,
		&updated.SKU,
		&updated.Price,
		&updated.CategoryID,
		&updated.Unit,
		&updated.IsActive,
		&deletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return nil, appErr.ErrProductAlreadyExists
		}
		return nil, err
	}

	if deletedAt.Valid {
		updated.DeletedAt = &deletedAt.Time
	}

	return &updated, nil
}

func (r *PostgresProductRepository) Delete(
	ctx context.Context,
	tenantID string,
	branchID string,
	productID string,
) error {
	query := `
		UPDATE products
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE tenant_id = $1 AND branch_id = $2 AND product_id = $3 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, tenantID, branchID, productID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return appErr.ErrProductNotFound
	}

	return nil
}