package handler

import (
	"context"
	"net/http"
	"pos-service/internal/domain"
	"time"

	"github.com/gin-gonic/gin"
)

type PosHandler struct {
	posService PosService
}

type PosService interface {
	GetHealth(ctx context.Context) error
	GetHealthByTenantID(ctx context.Context, tenantID string) error
	GetBranchesByTenantID(ctx context.Context, tenantID string) ([]domain.BranchResponse, error)
}

func NewPosHandler(s PosService) *PosHandler {
	return &PosHandler{
		posService: s,
	}
}

func (h *PosHandler) GetHealth(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
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

func (h *PosHandler) GetBranchesByTenantID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	tenantID := c.Param("tenant_id")
	if tenantID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "tenant_id is required",
		})
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

func (h *PosHandler) GetHealthByTenantID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	tenantID := c.Param("tenant_id")
	if tenantID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "tenant_id is required",
		})
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
