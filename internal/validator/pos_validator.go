package validator

import (
	"errors"
	"regexp"
)

type PosValidator struct{}

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