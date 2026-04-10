package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type PosHandler struct {
	posService PosService
	validator  TenantValidator
}

type PosService interface {
	GetHealth(ctx context.Context) error
	GetHealthByTenantID(ctx context.Context, tenantID string) error
}

type TenantValidator interface {
	TenantIDValidation(tenantID string) error
}

func NewPosHandler(s PosService, v TenantValidator) *PosHandler {
	return &PosHandler{
		posService: s,
		validator:  v,
	}
}

func (h *PosHandler) GetHealth(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	err := h.posService.GetHealth(ctx)
	if err != nil {
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
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	tenantID := c.Param("tenant_id")
	if tenantID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "tenant_id is required",
		})
		return
	}

	if err := h.validator.TenantIDValidation(tenantID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := h.posService.GetHealthByTenantID(ctx, tenantID)
	if err != nil {
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
	})
}