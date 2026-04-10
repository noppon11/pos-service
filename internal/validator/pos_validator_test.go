package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTenantIDValidation(t *testing.T) {
	v := &PosValidator{}

	tests := []struct {
		name      string
		input     string
		wantError bool
	}{
		{"valid dash", "aura-bkk", false},
		{"valid underscore", "tenant_001", false},
		{"empty", "", true},
		{"too short", "a", true},
		{"uppercase", "AURA", true},
		{"invalid char", "***", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.TenantIDValidation(tt.input)

			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}