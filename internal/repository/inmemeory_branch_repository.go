package repository

import (
	"context"

	"pos-service/internal/domain"
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
				},
				{
					BranchID:   "bkk-002",
					BranchName: "Aura Ari",
					Status:     "inactive",
				},
			},
			"aura-cnx": {
				{
					BranchID:   "cnx-001",
					BranchName: "Aura Chiang Mai",
					Status:     "active",
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

func (r *InMemoryBranchRepository) ByID(ctx context.Context, branchID string) (*domain.BranchResponse, error) {
	for _, branches := range r.data {
		for _, b := range branches {
			if b.BranchID == branchID {
				return &b, nil
			}
		}
	}
	return nil, nil
}