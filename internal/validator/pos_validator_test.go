package validator

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTenantIDValidation(t *testing.T) {
	v := &PosValidator{}

	tests := []struct {
		name        string
		input       string
		wantError   bool
		expectedErr string
	}{
		// ✅ valid cases
		{"valid dash", "aura-bkk", false, ""},
		{"valid underscore", "tenant_001", false, ""},
		{"valid mixed", "aura_bkk-01", false, ""},

		// ❗ boundary
		{"min length 3", "abc", false, ""},
		{"max length 50", strings.Repeat("a", 50), false, ""},

		// ❌ invalid cases
		{"empty", "", true, "tenant_id is required"},
		{"too short", "a", true, "tenant_id must be 3-50 chars, lowercase letters, numbers, underscore or dash only"},
		{"too long", strings.Repeat("a", 51), true, "tenant_id must be 3-50 chars, lowercase letters, numbers, underscore or dash only"},
		{"uppercase", "AURA", true, "tenant_id must be 3-50 chars, lowercase letters, numbers, underscore or dash only"},
		{"invalid char", "***", true, "tenant_id must be 3-50 chars, lowercase letters, numbers, underscore or dash only"},
		{"space inside", "aura bkk", true, "tenant_id must be 3-50 chars, lowercase letters, numbers, underscore or dash only"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.TenantIDValidation(tt.input)

			if tt.wantError {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}