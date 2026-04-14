package repository

import (
	"context"

	"pos-service/internal/domain"
	appErr "pos-service/internal/errors"
)

type InMemoryBranchRepository struct {
	data map[string][]domain.BranchResponse
}

func NewInMemoryBranchRepository() *InMemoryBranchRepository {
	return &InMemoryBranchRepository{
		data: map[string][]domain.BranchResponse{
			"aura-bkk": {
				{
					BranchID:   "bkk-001",
					BranchName: "Aura Siam",
					Status:     "active",
					Timezone:   "Asia/Bangkok",
					Currency:   "THB",    
				},
				{
					BranchID:   "bkk-002",
					BranchName: "Aura Ari",
					Status:     "inactive",
					Timezone:   "Asia/Bangkok",
					Currency:   "THB",    
				},
			},
			"aura-cnx": {
				{
					BranchID:   "cnx-001",
					BranchName: "Aura Chiang Mai",
					Status:     "active",
					Timezone:   "Asia/Chaing Mai",
					Currency:   "THB",    
				},
			},
		},
	}
}

func (r *InMemoryBranchRepository) ListByTenantID(ctx context.Context, tenantID string) ([]domain.BranchResponse, error) {
	branches, ok := r.data[tenantID]
	if !ok {
		return []domain.BranchResponse{}, nil
	}
	return branches, nil
}

func (r *InMemoryBranchRepository) GetByTenantIDAndBranchID(ctx context.Context, tenantID, branchID string) (*domain.BranchResponse, error) {
	branches, ok := r.data[tenantID]
	if !ok {
		return nil, appErr.ErrTenantNotFound
	}

	for i := range branches {
		if branches[i].BranchID == branchID {
			return &branches[i], nil
		}
	}

	return nil, appErr.ErrBranchNotFound
}