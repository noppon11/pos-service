package errors
import "errors"

var (
	// Infrastructure / Config
	ErrDBNotConfigured         = errors.New("database is not configured")
	ErrBranchRepoNotConfigured = errors.New("branch repository is not configured")
	ErrProductRepoNotConfigured = errors.New("product repository not configured")
	ErrValidatorNotConfigured  = errors.New("validator not configured")
	ErrProductIDRequired    = errors.New("product_id is required")
	ErrProductAlreadyExists = errors.New("product already exists")
	ErrProductNameRequired = errors.New("product name is required")
	ErrProductSKURequired  = errors.New("product sku is required")
	ErrInvalidProductPrice = errors.New("invalid product price")
	ErrCategoryIDRequired  = errors.New("category_id is required")
	ErrUnitRequired        = errors.New("unit is required")

	// Validation
	ErrInvalidBranchStatus           = errors.New("status must be active or inactive")
	ErrInvalidBranchCurrency         = errors.New("currency must be 3 uppercase letters")
	ErrInvalidBranchCurrencyRequired = errors.New("currency is required")
	ErrInvalidBranchTimezoneRequired = errors.New("timezone is required")
	ErrBranchIDRequired              = errors.New("branch_id is required")
	ErrBranchNameRequired            = errors.New("branch_name is required")
	ErrTenantIDRequired				= errors.New("tenant_id is required")
    ErrTenantNotFound 				 = errors.New("tenant not found")
    ErrBranchNotFound 				 = errors.New("branch not found")
	ErrProductNotFound         = errors.New("product not found")
	ErrCreateProductFailed     = errors.New("failed to create product")
)