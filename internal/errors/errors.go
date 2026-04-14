package errors
import "errors"

var (
	// Infrastructure / Config
	ErrDBNotConfigured         = errors.New("database is not configured")
	ErrBranchRepoNotConfigured = errors.New("branch repository is not configured")
	ErrValidatorNotConfigured  = errors.New("validator not configured")

	// Validation
	ErrInvalidBranchStatus           = errors.New("status must be active or inactive")
	ErrInvalidBranchCurrency         = errors.New("currency must be 3 uppercase letters")
	ErrInvalidBranchCurrencyRequired = errors.New("currency is required")
	ErrInvalidBranchTimezoneRequired = errors.New("timezone is required")
	ErrBranchIDRequired              = errors.New("branch_id is required")
	ErrBranchNameRequired            = errors.New("branch_name is required")
    ErrTenantNotFound 				 = errors.New("tenant not found")
    ErrBranchNotFound 				 = errors.New("branch not found")
)