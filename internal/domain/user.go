package domain

import "time"

type UserRole string

const (
	RoleOwner   UserRole = "owner"
	RoleManager UserRole = "manager"
	RoleStaff   UserRole = "staff"
)

type User struct {
	ID           string
	Email        string
	PasswordHash string
	FullName     string
	TenantID     string
	Role         UserRole
	IsActive     bool
	BranchIDs    []string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}