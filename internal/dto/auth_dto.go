package dto

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthUserResponse struct {
	ID        string   `json:"id"`
	Email     string   `json:"email"`
	FullName  string   `json:"full_name"`
	TenantID  string   `json:"tenant_id"`
	BranchIDs []string `json:"branch_ids"`
	Role      string   `json:"role"`
}

type LoginResponse struct {
	AccessToken string           `json:"access_token"`
	TokenType   string           `json:"token_type"`
	ExpiresIn   int64            `json:"expires_in"`
	User        AuthUserResponse `json:"user"`
}

type MeResponse struct {
	ID        string   `json:"id"`
	Email     string   `json:"email"`
	FullName  string   `json:"full_name"`
	TenantID  string   `json:"tenant_id"`
	BranchIDs []string `json:"branch_ids"`
	Role      string   `json:"role"`
}