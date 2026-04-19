package validator

import (
	"errors"
	"regexp"
	"strings"

	"pos-service/internal/domain"
	appErr "pos-service/internal/errors"
)

type PosValidator struct{}

type TenantValidator interface {
	TenantIDValidation(tenantID string) error
	BranchIDValidation(branchID string) error
	ProductIDValidation(productID string) error
}

type BranchValidator interface {
	ValidateBranch(branch domain.BranchResponse) error
}

type ProductValidator interface {
	ValidateProduct(product domain.Product) error
}

var (
	IDRegex       = regexp.MustCompile(`^[a-z0-9_-]{3,50}$`)
	CurrencyRegex = regexp.MustCompile(`^[A-Z]{3}$`)
)

func (v *PosValidator) TenantIDValidation(tenantID string) error {
	if strings.TrimSpace(tenantID) == "" {
		return errors.New("tenant_id is required")
	}

	if !IDRegex.MatchString(tenantID) {
		return errors.New("tenant_id must be 3-50 chars, lowercase letters, numbers, underscore or dash only")
	}

	return nil
}

func (v *PosValidator) BranchIDValidation(branchID string) error {
	if strings.TrimSpace(branchID) == "" {
		return errors.New("branch_id is required")
	}

	if !IDRegex.MatchString(branchID) {
		return errors.New("branch_id must be 3-50 chars, lowercase letters, numbers, underscore or dash only")
	}

	return nil
}

func (v *PosValidator) ProductIDValidation(productID string) error {
	if strings.TrimSpace(productID) == "" {
		return errors.New("product_id is required")
	}

	if !IDRegex.MatchString(productID) {
		return errors.New("product_id must be 3-50 chars, lowercase letters, numbers, underscore or dash only")
	}

	return nil
}

func (v *PosValidator) ValidateBranch(branch domain.BranchResponse) error {
	if strings.TrimSpace(branch.BranchID) == "" {
		return appErr.ErrBranchIDRequired
	}

	if strings.TrimSpace(branch.BranchName) == "" {
		return appErr.ErrBranchNameRequired
	}

	if !isValidBranchStatus(branch.Status) {
		return appErr.ErrInvalidBranchStatus
	}

	if strings.TrimSpace(branch.Timezone) == "" {
		return appErr.ErrInvalidBranchTimezoneRequired
	}

	if strings.TrimSpace(branch.Currency) == "" {
		return appErr.ErrInvalidBranchCurrencyRequired
	}

	if !CurrencyRegex.MatchString(branch.Currency) {
		return appErr.ErrInvalidBranchCurrency
	}

	return nil
}

func (v *PosValidator) ValidateProduct(product domain.Product) error {
     if strings.TrimSpace(product.Name) == "" {
         return appErr.ErrProductNameRequired
     }
     if strings.TrimSpace(product.SKU) == "" {
         return appErr.ErrProductSKURequired
     }
    if product.Price <= 0 {
         return appErr.ErrInvalidProductPrice
     }
     if strings.TrimSpace(product.CategoryID) == "" {
         return appErr.ErrCategoryIDRequired
     }
     if strings.TrimSpace(product.Unit) == "" {
         return appErr.ErrUnitRequired
     }
     return nil
 }

func isValidBranchStatus(status string) bool {
	switch status {
	case "active", "inactive":
		return true
	default:
		return false
	}
}