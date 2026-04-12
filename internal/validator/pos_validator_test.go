package validator

import (
	"strings"
	"testing"

	"pos-service/internal/domain"

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
		{"valid mixed", "aura_bkk-001", false},
		{"min length 3", "abc", false},
		{"max length 50", strings.Repeat("a", 50), false},
		{"empty", "", true},
		{"too short", "a", true},
		{"too long", strings.Repeat("a", 51), true},
		{"uppercase", "AURA", true},
		{"invalid char", "***", true},
		{"space inside", "aura bkk", true},
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

func TestBranchValidation(t *testing.T) {
	v := &PosValidator{}

	tests := []struct {
		name      string
		input     domain.BranchResponse
		wantError bool
		wantMsg   string
	}{
		{
			name: "valid active",
			input: domain.BranchResponse{
				BranchID:   "bkk-001",
				BranchName: "Aura Siam",
				Status:     "active",
			},
			wantError: false,
		},
		{
			name: "valid inactive",
			input: domain.BranchResponse{
				BranchID:   "bkk-002",
				BranchName: "Aura Ari",
				Status:     "inactive",
			},
			wantError: false,
		},
		{
			name: "empty branch_id",
			input: domain.BranchResponse{
				BranchID:   "",
				BranchName: "Aura Siam",
				Status:     "active",
			},
			wantError: true,
			wantMsg:   "branch_id is required",
		},
		{
			name: "branch_id only spaces",
			input: domain.BranchResponse{
				BranchID:   "   ",
				BranchName: "Aura Siam",
				Status:     "active",
			},
			wantError: true,
			wantMsg:   "branch_id is required",
		},
		{
			name: "empty branch_name",
			input: domain.BranchResponse{
				BranchID:   "bkk-001",
				BranchName: "",
				Status:     "active",
			},
			wantError: true,
			wantMsg:   "branch_name is required",
		},
		{
			name: "branch_name only spaces",
			input: domain.BranchResponse{
				BranchID:   "bkk-001",
				BranchName: "   ",
				Status:     "active",
			},
			wantError: true,
			wantMsg:   "branch_name is required",
		},
		{
			name: "invalid status",
			input: domain.BranchResponse{
				BranchID:   "bkk-001",
				BranchName: "Aura Siam",
				Status:     "pending",
			},
			wantError: true,
			wantMsg:   "status must be active or inactive",
		},
		{
			name: "empty status",
			input: domain.BranchResponse{
				BranchID:   "bkk-001",
				BranchName: "Aura Siam",
				Status:     "",
			},
			wantError: true,
			wantMsg:   "status must be active or inactive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.BranchValidation(tt.input)

			if tt.wantError {
				assert.Error(t, err)
				assert.Equal(t, tt.wantMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}