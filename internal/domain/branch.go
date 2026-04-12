package domain

type BranchResponse struct {
	BranchID   string `json:"branch_id"`
	BranchName string `json:"branch_name"`
	Status     string `json:"status"`
}

type ListBranchesResponse struct {
	TenantID string   `json:"tenant_id"`
	Data     []BranchResponse `json:"data"`
}