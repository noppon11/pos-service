package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"pos-service/internal/domain"
	"pos-service/internal/service"
)

type TokenParser interface {
	ParseToken(token string) (*domain.AuthClaims, error)
}

func RequireAuth(tokenParser TokenParser) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing authorization header",
			})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization header",
			})
			return
		}

		claims, err := tokenParser.ParseToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token",
			})
			return
		}

		c.Set("auth_claims", claims)
		c.Next()
	}
}

func RequireTenantBranchAccess(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		raw, exists := c.Get("auth_claims")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		}

		claims, ok := raw.(*domain.AuthClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		}

		tenantID := c.Param("tenantID")
		branchID := c.Param("branchID")

		if err := authService.AuthorizeTenantBranch(claims, tenantID, branchID); err != nil {
			switch err {
			case service.ErrForbiddenTenant, service.ErrForbiddenBranch:
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"error": "forbidden",
				})
				return
			default:
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "internal server error",
				})
				return
			}
		}

		c.Next()
	}
}