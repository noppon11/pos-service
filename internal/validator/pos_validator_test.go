package validator

import (
	"testing"

	"pos-service/internal/domain"

	"github.com/stretchr/testify/assert"
)

func TestTenantIDValidation(t *testing.T) {
	v := &PosValidator{}

	tests := []struct {
		name     string
		tenantID string
		wantErr  bool
	}{
		{"empty", "", true},
		{"invalid uppercase", "AURA", true},
		{"invalid symbol", "aura@bkk", true},
		{"valid", "aura-bkk", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.TenantIDValidation(tt.tenantID)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestBranchIDValidation(t *testing.T) {
	v := &PosValidator{}

	tests := []struct {
		name     string
		branchID string
		wantErr  bool
	}{
		{"empty", "", true},
		{"invalid uppercase", "BKK-001", true},
		{"invalid symbol", "bkk@001", true},
		{"valid", "bkk-001", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.BranchIDValidation(tt.branchID)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestProductIDValidation(t *testing.T) {
	v := &PosValidator{}

	tests := []struct {
		name      string
		productID string
		wantErr   bool
	}{
		{"empty", "", true},
		{"invalid uppercase", "PROD-001", true},
		{"invalid symbol", "prod@001", true},
		{"valid", "prod-001", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.ProductIDValidation(tt.productID)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestValidateBranch(t *testing.T) {
	v := &PosValidator{}

	tests := []struct {
		name    string
		branch  domain.BranchResponse
		wantErr bool
	}{
		{
			name: "missing branch id",
			branch: domain.BranchResponse{
				BranchName: "Aura Siam",
				Status:     "active",
				Timezone:   "Asia/Bangkok",
				Currency:   "THB",
			},
			wantErr: true,
		},
		{
			name: "missing branch name",
			branch: domain.BranchResponse{
				BranchID:  "bkk-001",
				Status:    "active",
				Timezone:  "Asia/Bangkok",
				Currency:  "THB",
			},
			wantErr: true,
		},
		{
			name: "invalid status",
			branch: domain.BranchResponse{
				BranchID:   "bkk-001",
				BranchName: "Aura Siam",
				Status:     "paused",
				Timezone:   "Asia/Bangkok",
				Currency:   "THB",
			},
			wantErr: true,
		},
		{
			name: "missing timezone",
			branch: domain.BranchResponse{
				BranchID:   "bkk-001",
				BranchName: "Aura Siam",
				Status:     "active",
				Currency:   "THB",
			},
			wantErr: true,
		},
		{
			name: "invalid currency",
			branch: domain.BranchResponse{
				BranchID:   "bkk-001",
				BranchName: "Aura Siam",
				Status:     "active",
				Timezone:   "Asia/Bangkok",
				Currency:   "thb",
			},
			wantErr: true,
		},
		{
			name: "success",
			branch: domain.BranchResponse{
				BranchID:   "bkk-001",
				BranchName: "Aura Siam",
				Status:     "active",
				Timezone:   "Asia/Bangkok",
				Currency:   "THB",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.ValidateBranch(tt.branch)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestValidateProduct(t *testing.T) {
	v := &PosValidator{}

	tests := []struct {
		name    string
		product domain.ProductResponse
		wantErr bool
	}{
		{
			name: "missing name",
			product: domain.ProductResponse{
				SKU:        "BOT-50",
				Price:      3500,
				CategoryID: "treatment",
				Unit:       "unit",
				IsActive:   true,
			},
			wantErr: true,
		},
		{
			name: "missing sku",
			product: domain.ProductResponse{
				Name:       "Botox 50u",
				Price:      3500,
				CategoryID: "treatment",
				Unit:       "unit",
				IsActive:   true,
			},
			wantErr: true,
		},
		{
			name: "negative price",
			product: domain.ProductResponse{
				Name:       "Botox 50u",
				SKU:        "BOT-50",
				Price:      -1,
				CategoryID: "treatment",
				Unit:       "unit",
				IsActive:   true,
			},
			wantErr: true,
		},
		{
			name: "missing category id",
			product: domain.ProductResponse{
				Name:     "Botox 50u",
				SKU:      "BOT-50",
				Price:    3500,
				Unit:     "unit",
				IsActive: true,
			},
			wantErr: true,
		},
		{
			name: "missing unit",
			product: domain.ProductResponse{
				Name:       "Botox 50u",
				SKU:        "BOT-50",
				Price:      3500,
				CategoryID: "treatment",
				IsActive:   true,
			},
			wantErr: true,
		},
		{
			name: "success",
			product: domain.ProductResponse{
				ProductID:  "prod-001",
				Name:       "Botox 50u",
				SKU:        "BOT-50",
				Price:      3500,
				CategoryID: "treatment",
				Unit:       "unit",
				IsActive:   true,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.ValidateProduct(tt.product)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}