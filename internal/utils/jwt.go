package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"pos-service/internal/domain"
)

type JWTManager struct {
	secretKey []byte
	expiresIn time.Duration
}

func NewJWTManager(secret string, expiresIn time.Duration) *JWTManager {
	return &JWTManager{
		secretKey: []byte(secret),
		expiresIn: expiresIn,
	}
}

type CustomClaims struct {
	UserID    string   `json:"user_id"`
	Email     string   `json:"email"`
	FullName  string   `json:"full_name"`
	TenantID  string   `json:"tenant_id"`
	BranchIDs []string `json:"branch_ids"`
	Role      string   `json:"role"`
	jwt.RegisteredClaims
}

func (m *JWTManager) GenerateToken(user *domain.User) (string, int64, error) {
	expireAt := time.Now().Add(m.expiresIn)

	claims := CustomClaims{
		UserID:    user.ID,
		Email:     user.Email,
		FullName:  user.FullName,
		TenantID:  user.TenantID,
		BranchIDs: user.BranchIDs,
		Role:      string(user.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.ID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(m.secretKey)
	if err != nil {
		return "", 0, err
	}

	return signed, int64(m.expiresIn.Seconds()), nil
}

func (m *JWTManager) ParseToken(tokenString string) (*domain.AuthClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return m.secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	var exp int64
	if claims.ExpiresAt != nil {
		exp = claims.ExpiresAt.Unix()
	}

	return &domain.AuthClaims{
		UserID:    claims.UserID,
		Email:     claims.Email,
		FullName:  claims.FullName,
		TenantID:  claims.TenantID,
		BranchIDs: claims.BranchIDs,
		Role:      claims.Role,
		ExpiresAt: exp,
	}, nil
}