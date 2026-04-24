package domain

type AuthClaims struct {
	UserID    string   `json:"user_id"`
	Email     string   `json:"email"`
	FullName  string   `json:"full_name"`
	TenantID  string   `json:"tenant_id"`
	BranchIDs []string `json:"branch_ids"`
	Role      string   `json:"role"`
	ExpiresAt int64    `json:"expires_at"`
}