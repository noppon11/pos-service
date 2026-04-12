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
}

type BranchValidator interface {
	BranchValidation(branch domain.BranchResponse) error
}

var tenantIDRegex = regexp.MustCompile(`^[a-z0-9_-]{3,50}$`)

func (v *PosValidator) TenantIDValidation(tenantID string) error {
	if tenantID == "" {
		return errors.New("tenant_id is required")
	}

	if !tenantIDRegex.MatchString(tenantID) {
		return errors.New("tenant_id must be 3-50 chars, lowercase letters, numbers, underscore or dash only")
	}

	return nil
}

func (v *PosValidator) BranchValidation(branch domain.BranchResponse) error {
	if strings.TrimSpace(branch.BranchID) == "" {
		return errors.New("branch_id is required")
	}

	if strings.TrimSpace(branch.BranchName) == "" {
		return errors.New("branch_name is required")
	}

	if !isValidBranchStatus(branch.Status) {
		return errors.New("status must be active or inactive")
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