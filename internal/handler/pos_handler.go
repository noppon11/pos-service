package handler

import (
	"context"
	"net/http"
	"pos-service/internal/domain"
	"time"

	"github.com/gin-gonic/gin"
)

const requestTimeout = 2 * time.Second

type PosHandler struct {
	posService PosService
	validator  TenantValidator
}

type TenantValidator interface {
	TenantIDValidation(tenantID string) error
}

type PosService interface {
	GetHealth(ctx context.Context) error
	GetHealthByTenantID(ctx context.Context, tenantID string) error
	GetBranchesByTenantID(ctx context.Context, tenantID string) ([]domain.BranchResponse, error)
}

func NewPosHandler(s PosService, v TenantValidator) *PosHandler {
	return &PosHandler{
		posService: s,
		validator:  v,
	}
}

func (h *PosHandler) GetHealth(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	if err := h.posService.GetHealth(ctx); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"service":   "pos-service",
			"status":    "unhealthy",
			"timestamp": time.Now().Unix(),
			"error":     err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"service":   "pos-service",
		"status":    "ok",
		"timestamp": time.Now().Unix(),
	})
}

func (h *PosHandler) GetHealthByTenantID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	tenantID, ok := h.getValidTenantID(c)
	if !ok {
		return
	}

	if err := h.posService.GetHealthByTenantID(ctx, tenantID); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"service":   "pos-service",
			"status":    "unhealthy",
			"tenant_id": tenantID,
			"timestamp": time.Now().Unix(),
			"error":     err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"service":   "pos-service",
		"status":    "ok",
		"tenant_id": tenantID,
		"timestamp": time.Now().Unix(),
	})
}

func (h *PosHandler) GetBranchesByTenantID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	tenantID, ok := h.getValidTenantID(c)
	if !ok {
		return
	}

	data, err := h.posService.GetBranchesByTenantID(ctx, tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, domain.ListBranchesResponse{
		TenantID: tenantID,
		Data:     data,
	})
}

func (h *PosHandler) getValidTenantID(c *gin.Context) (string, bool) {
	tenantID := c.Param("tenant_id")
	if tenantID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "tenant_id is required",
		})
		return "", false
	}

	if err := h.validator.TenantIDValidation(tenantID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return "", false
	}

	return tenantID, true
}