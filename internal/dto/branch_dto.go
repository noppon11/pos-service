package dto
import "pos-service/internal/domain"

type BranchResponseDTO struct {
	BranchID   string `json:"branch_id"`
	BranchName string `json:"branch_name"`
	Status     string `json:"status"`
	Timezone   string `json:"timezone"`
	Currency   string `json:"currency"`
}

type ListBranchesResponseDTO struct {
	TenantID string               `json:"tenant_id"`
	Data     []BranchResponseDTO  `json:"data"`
}

func ToBranchResponseDTO(b domain.BranchResponse) BranchResponseDTO {
	return BranchResponseDTO{
		BranchID:   b.BranchID,
		BranchName: b.BranchName,
		Status:     b.Status,
		Timezone:   b.Timezone,
		Currency:   b.Currency,
	}
}

func ToListBranchesResponseDTO(tenantID string, branches []domain.BranchResponse) ListBranchesResponseDTO {
	data := make([]BranchResponseDTO, 0, len(branches))

	for _, b := range branches {
		data = append(data, ToBranchResponseDTO(b))
	}

	return ListBranchesResponseDTO{
		TenantID: tenantID,
		Data:     data,
	}
}