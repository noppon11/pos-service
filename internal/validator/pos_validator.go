package validator

import (
	"errors"
	"regexp"
	"strings"

	"pos-service/internal/domain"
)

type PosValidator struct{}

type TenantValidator interface {
	TenantIDValidation(tenantID string) error
	BranchIDValidation(branchID string) error
}

type BranchValidator interface {
	BranchValidation(branch domain.BranchResponse) error
}

var (
	ErrBranchIDRequired              = errors.New("branch_id is required")
	ErrBranchNameRequired            = errors.New("branch_name is required")
	ErrInvalidBranchStatus           = errors.New("status must be active or inactive")
	ErrInvalidBranchCurrency         = errors.New("currency must be 3 uppercase letters")
	ErrInvalidBranchCurrencyRequired = errors.New("currency is required")
	ErrInvalidBranchTimezoneRequired = errors.New("timezone is required")
)

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

func (v *PosValidator) BranchValidation(branch domain.BranchResponse) error {
	if strings.TrimSpace(branch.BranchID) == "" {
		return ErrBranchIDRequired
	}

	if strings.TrimSpace(branch.BranchName) == "" {
		return ErrBranchNameRequired
	}

	if !isValidBranchStatus(branch.Status) {
		return ErrInvalidBranchStatus
	}

	if strings.TrimSpace(branch.Timezone) == "" {
		return ErrInvalidBranchTimezoneRequired
	}

	if strings.TrimSpace(branch.Currency) == "" {
		return ErrInvalidBranchCurrencyRequired
	}

	if !CurrencyRegex.MatchString(branch.Currency) {
		return ErrInvalidBranchCurrency 
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