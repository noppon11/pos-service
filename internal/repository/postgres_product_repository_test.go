package repository

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"pos-service/internal/domain"
	appErr "pos-service/internal/errors"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestPostgresProductRepository_Create_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewPostgresProductRepository(db)

	product := domain.Product{
		Name:       "Botox 100u",
		SKU:        "BOT-100",
		Price:      6500,
		CategoryID: "treatment",
		Unit:       "unit",
		IsActive:   true,
	}

	rows := sqlmock.NewRows([]string{
		"product_id", "name", "sku", "price", "category_id", "unit", "is_active",
	}).AddRow(
		"generated-uuid",
		"Botox 100u",
		"BOT-100",
		6500,
		"treatment",
		"unit",
		true,
	)

	mock.ExpectQuery(regexp.QuoteMeta(`
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
	`)).
		WithArgs(
			sqlmock.AnyArg(),
			"aura-bkk",
			"bkk-001",
			product.Name,
			product.SKU,
			product.Price,
			product.CategoryID,
			product.Unit,
			product.IsActive,
		).
		WillReturnRows(rows)

	created, err := repo.Create(context.Background(), "aura-bkk", "bkk-001", product)
	assert.NoError(t, err)
	assert.NotNil(t, created)
	assert.Equal(t, "generated-uuid", created.ProductID)
	assert.Equal(t, "Botox 100u", created.Name)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresProductRepository_ListByTenantIDAndBranchID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewPostgresProductRepository(db)

	filter := ProductListFilter{
		Page:  1,
		Limit: 20,
	}

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT COUNT(*)
		FROM products
		WHERE tenant_id = $1
		  AND branch_id = $2
		  AND deleted_at IS NULL
		  AND ($3 = '' OR category_id = $3)
	`)).
		WithArgs("aura-bkk", "bkk-001", "").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

	rows := sqlmock.NewRows([]string{
		"product_id", "name", "sku", "price", "category_id", "unit", "is_active", "deleted_at",
	}).AddRow(
		"prod-001", "Botox 50u", "BOT-50", 3500, "treatment", "unit", true, nil,
	).AddRow(
		"prod-002", "Filler 1cc", "FIL-1", 12000, "treatment", "unit", true, nil,
	)

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT product_id, name, sku, price, category_id, unit, is_active, deleted_at
		FROM products
		WHERE tenant_id = $1
		  AND branch_id = $2
		  AND deleted_at IS NULL
		  AND ($3 = '' OR category_id = $3)
		ORDER BY created_at DESC
		LIMIT $4 OFFSET $5
	`)).
		WithArgs("aura-bkk", "bkk-001", "", 20, 0).
		WillReturnRows(rows)

	products, total, err := repo.ListByTenantIDAndBranchID(
		context.Background(),
		"aura-bkk",
		"bkk-001",
		filter,
	)

	assert.NoError(t, err)
	assert.Equal(t, 2, total)
	assert.Len(t, products, 2)
	assert.Equal(t, "prod-001", products[0].ProductID)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresProductRepository_GetByTenantIDBranchIDAndProductID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewPostgresProductRepository(db)

	rows := sqlmock.NewRows([]string{
		"product_id", "name", "sku", "price", "category_id", "unit", "is_active", "deleted_at",
	}).AddRow(
		"prod-001", "Botox 50u", "BOT-50", 3500, "treatment", "unit", true, nil,
	)

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT product_id, name, sku, price, category_id, unit, is_active, deleted_at
		FROM products
		WHERE tenant_id = $1
		  AND branch_id = $2
		  AND product_id = $3
		  AND deleted_at IS NULL
	`)).
		WithArgs("aura-bkk", "bkk-001", "prod-001").
		WillReturnRows(rows)

	product, err := repo.GetByTenantIDBranchIDAndProductID(context.Background(), "aura-bkk", "bkk-001", "prod-001")
	assert.NoError(t, err)
	assert.NotNil(t, product)
	assert.Equal(t, "prod-001", product.ProductID)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresProductRepository_GetByTenantIDBranchIDAndProductID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewPostgresProductRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT product_id, name, sku, price, category_id, unit, is_active, deleted_at
		FROM products
		WHERE tenant_id = $1
		  AND branch_id = $2
		  AND product_id = $3
		  AND deleted_at IS NULL
	`)).
		WithArgs("aura-bkk", "bkk-001", "prod-001").
		WillReturnError(sql.ErrNoRows)

	product, err := repo.GetByTenantIDBranchIDAndProductID(context.Background(), "aura-bkk", "bkk-001", "prod-001")
	assert.NoError(t, err)
	assert.Nil(t, product)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresProductRepository_Update_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewPostgresProductRepository(db)

	product := domain.Product{
		Name:       "Botox 120u",
		SKU:        "BOT-120",
		Price:      7500,
		CategoryID: "treatment",
		Unit:       "unit",
		IsActive:   true,
	}

	rows := sqlmock.NewRows([]string{
		"product_id", "name", "sku", "price", "category_id", "unit", "is_active", "deleted_at",
	}).AddRow(
		"prod-001", "Botox 120u", "BOT-120", 7500, "treatment", "unit", true, nil,
	)

	mock.ExpectQuery(regexp.QuoteMeta(`
		UPDATE products
		SET
			name = $4,
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
	`)).
		WithArgs(
			"aura-bkk",
			"bkk-001",
			"prod-001",
			product.Name,
			product.SKU,
			product.Price,
			product.CategoryID,
			product.Unit,
			product.IsActive,
		).
		WillReturnRows(rows)

	updated, err := repo.Update(context.Background(), "aura-bkk", "bkk-001", "prod-001", product)
	assert.NoError(t, err)
	assert.NotNil(t, updated)
	assert.Equal(t, "prod-001", updated.ProductID)
	assert.Equal(t, "Botox 120u", updated.Name)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresProductRepository_Update_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewPostgresProductRepository(db)

	product := domain.Product{
		Name:       "Botox 120u",
		SKU:        "BOT-120",
		Price:      7500,
		CategoryID: "treatment",
		Unit:       "unit",
		IsActive:   true,
	}

	mock.ExpectQuery(regexp.QuoteMeta(`
		UPDATE products
		SET
			name = $4,
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
	`)).
		WithArgs(
			"aura-bkk",
			"bkk-001",
			"prod-001",
			product.Name,
			product.SKU,
			product.Price,
			product.CategoryID,
			product.Unit,
			product.IsActive,
		).
		WillReturnError(sql.ErrNoRows)

	updated, err := repo.Update(context.Background(), "aura-bkk", "bkk-001", "prod-001", product)
	assert.NoError(t, err)
	assert.Nil(t, updated)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresProductRepository_Delete_SoftDeleteSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewPostgresProductRepository(db)

	mock.ExpectExec(regexp.QuoteMeta(`
		UPDATE products
		SET deleted_at = NOW(),
		    updated_at = NOW()
		WHERE tenant_id = $1
		  AND branch_id = $2
		  AND product_id = $3
		  AND deleted_at IS NULL
	`)).
		WithArgs("aura-bkk", "bkk-001", "prod-001").
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Delete(context.Background(), "aura-bkk", "bkk-001", "prod-001")
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresProductRepository_Delete_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewPostgresProductRepository(db)

	mock.ExpectExec(regexp.QuoteMeta(`
		UPDATE products
		SET deleted_at = NOW(),
		    updated_at = NOW()
		WHERE tenant_id = $1
		  AND branch_id = $2
		  AND product_id = $3
		  AND deleted_at IS NULL
	`)).
		WithArgs("aura-bkk", "bkk-001", "prod-001").
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = repo.Delete(context.Background(), "aura-bkk", "bkk-001", "prod-001")
	assert.ErrorIs(t, err, appErr.ErrProductNotFound)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresProductRepository_ScanDeletedAt(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewPostgresProductRepository(db)
	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"product_id", "name", "sku", "price", "category_id", "unit", "is_active", "deleted_at",
	}).AddRow(
		"prod-001", "Botox 50u", "BOT-50", 3500, "treatment", "unit", true, now,
	)

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT product_id, name, sku, price, category_id, unit, is_active, deleted_at
		FROM products
		WHERE tenant_id = $1
		  AND branch_id = $2
		  AND product_id = $3
		  AND deleted_at IS NULL
	`)).
		WithArgs("aura-bkk", "bkk-001", "prod-001").
		WillReturnRows(rows)

	product, err := repo.GetByTenantIDBranchIDAndProductID(context.Background(), "aura-bkk", "bkk-001", "prod-001")
	assert.NoError(t, err)
	assert.NotNil(t, product)
	assert.NotNil(t, product.DeletedAt)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresProductRepository_Create_DuplicateSKU(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewPostgresProductRepository(db)

	product := domain.Product{
		Name:       "Botox 100u",
		SKU:        "BOT-100",
		Price:      6500,
		CategoryID: "treatment",
		Unit:       "unit",
		IsActive:   true,
	}

	mock.ExpectQuery(regexp.QuoteMeta(`
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
	`)).
		WithArgs(
			sqlmock.AnyArg(),
			"aura-bkk",
			"bkk-001",
			product.Name,
			product.SKU,
			product.Price,
			product.CategoryID,
			product.Unit,
			product.IsActive,
		).
		WillReturnError(&pq.Error{Code: "23505"})

	created, err := repo.Create(context.Background(), "aura-bkk", "bkk-001", product)

	assert.Nil(t, created)
	assert.ErrorIs(t, err, appErr.ErrProductAlreadyExists)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresProductRepository_ListByTenantIDAndBranchID_WithCategoryFilterAndPagination(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewPostgresProductRepository(db)

	filter := ProductListFilter{
		Page:       2,
		Limit:      10,
		CategoryID: "treatment",
	}

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT COUNT(*)
		FROM products
		WHERE tenant_id = $1
		  AND branch_id = $2
		  AND deleted_at IS NULL
		  AND ($3 = '' OR category_id = $3)
	`)).
		WithArgs("aura-bkk", "bkk-001", "treatment").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(21))

	rows := sqlmock.NewRows([]string{
		"product_id", "name", "sku", "price", "category_id", "unit", "is_active", "deleted_at",
	}).AddRow(
		"prod-011", "Botox 50u", "BOT-50", 3500, "treatment", "unit", true, nil,
	)

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT product_id, name, sku, price, category_id, unit, is_active, deleted_at
		FROM products
		WHERE tenant_id = $1
		  AND branch_id = $2
		  AND deleted_at IS NULL
		  AND ($3 = '' OR category_id = $3)
		ORDER BY created_at DESC
		LIMIT $4 OFFSET $5
	`)).
		WithArgs("aura-bkk", "bkk-001", "treatment", 10, 10).
		WillReturnRows(rows)

	products, total, err := repo.ListByTenantIDAndBranchID(context.Background(), "aura-bkk", "bkk-001", filter)

	assert.NoError(t, err)
	assert.Len(t, products, 1)
	assert.Equal(t, 21, total)
	assert.Equal(t, "prod-011", products[0].ProductID)
	assert.NoError(t, mock.ExpectationsWereMet())
}