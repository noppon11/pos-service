package service

import (
	"context"
	"errors"
	"strings"

	"pos-service/internal/domain"
	"pos-service/internal/dto"
	"pos-service/internal/repository"
	"pos-service/internal/utils"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInactiveUser       = errors.New("inactive user")
	ErrForbiddenTenant    = errors.New("forbidden tenant")
	ErrForbiddenBranch    = errors.New("forbidden branch")
)

type TokenManager interface {
	GenerateToken(user *domain.User) (string, int64, error)
	ParseToken(token string) (*domain.AuthClaims, error)
}

type AuthService struct {
	userRepo     repository.UserRepository
	tokenManager TokenManager
}

func NewAuthService(
	userRepo repository.UserRepository,
	tokenManager TokenManager,
) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		tokenManager: tokenManager,
	}
}

func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.Password = strings.TrimSpace(req.Password)

	if req.Email == "" || req.Password == "" {
		return nil, ErrInvalidCredentials
	}

	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	if !user.IsActive {
		return nil, ErrInactiveUser
	}

	if err := utils.CheckPassword(user.PasswordHash, req.Password); err != nil {
		return nil, ErrInvalidCredentials
	}

	token, expiresIn, err := s.tokenManager.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   expiresIn,
		User: dto.AuthUserResponse{
			ID:        user.ID,
			Email:     user.Email,
			FullName:  user.FullName,
			TenantID:  user.TenantID,
			BranchIDs: user.BranchIDs,
			Role:      string(user.Role),
		},
	}, nil
}

func (s *AuthService) GetMe(ctx context.Context, claims *domain.AuthClaims) (*dto.MeResponse, error) {
	return &dto.MeResponse{
		ID:        claims.UserID,
		Email:     claims.Email,
		FullName:  claims.FullName,
		TenantID:  claims.TenantID,
		BranchIDs: claims.BranchIDs,
		Role:      claims.Role,
	}, nil
}

func (s *AuthService) AuthorizeTenantBranch(claims *domain.AuthClaims, tenantID, branchID string) error {
	if claims.TenantID != tenantID {
		return ErrForbiddenTenant
	}

	if claims.Role == string(domain.RoleOwner) {
		return nil
	}

	for _, b := range claims.BranchIDs {
		if b == branchID {
			return nil
		}
	}

	return ErrForbiddenBranch
}