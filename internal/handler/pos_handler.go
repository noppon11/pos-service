package handler

import (
	"context"
	"net/http"
	"pos-service/internal/dto"
	"pos-service/internal/domain"
	"time"

	"github.com/gin-gonic/gin"
)

const requestTimeout = 2 * time.Second

type PosHandler struct {
	posService PosService
	validator  Validator
}

type Validator interface {
	TenantIDValidation(tenantID string) error
	BranchIDValidation(branchID string) error
	ProductIDValidation(productID string) error
}

type PosService interface {
	GetHealth(ctx context.Context) error
	GetHealthByTenantID(ctx context.Context, tenantID string) error
	GetBranchDetail(ctx context.Context, tenantID string, branchID string) (*domain.BranchResponse, error)
	GetBranchesByTenantID(ctx context.Context, tenantID string) ([]domain.BranchResponse, error)

	GetProducts(ctx context.Context, tenantID string, branchID string) ([]domain.ProductResponse, error)
	GetProductByID(ctx context.Context, tenantID string, branchID string, productID string) (*domain.ProductResponse, error)
	CreateNewProduct(ctx context.Context, tenantID string, branchID string, req dto.CreateProductRequest) (*domain.ProductResponse, error)
	UpdateProduct(ctx context.Context, tenantID string, branchID string, productID string, req dto.UpdateProductRequest) (*domain.ProductResponse, error)
	DeleteProduct(ctx context.Context, tenantID string, branchID string, productID string) error
}

func NewPosHandler(s PosService, v Validator) *PosHandler {
	return &PosHandler{
		posService: s,
		validator:  v,
	}
}

func (h *PosHandler) newTimeoutContext(c *gin.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(c.Request.Context(), requestTimeout)
}

func (h *PosHandler) respondError(c *gin.Context, status int, err error) {
	c.JSON(status, gin.H{
		"error": err.Error(),
	})
}

func (h *PosHandler) respondValidationError(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, gin.H{
		"error": message,
	})
}

func toProductResponse(data *domain.ProductResponse) *domain.ProductResponse {
	if data == nil {
		return nil
	}

	return &domain.ProductResponse{
		ProductID:  data.ProductID,
		Name:       data.Name,
		SKU:        data.SKU,
		Price:      data.Price,
		CategoryID: data.CategoryID,
		Unit:       data.Unit,
		IsActive:   data.IsActive,
		DeletedAt:   data.DeletedAt,
	}
}

func (h *PosHandler) GetHealth(c *gin.Context) {
	ctx, cancel := h.newTimeoutContext(c)
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
	ctx, cancel := h.newTimeoutContext(c)
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
	ctx, cancel := h.newTimeoutContext(c)
	defer cancel()

	tenantID, ok := h.getValidTenantID(c)
	if !ok {
		return
	}

	branches, err := h.posService.GetBranchesByTenantID(ctx, tenantID)
	if err != nil {
		h.respondError(c, http.StatusInternalServerError, err)
		return
	}

	resp := dto.ToListBranchesResponseDTO(tenantID, branches)
	c.JSON(http.StatusOK, resp)
}

func (h *PosHandler) GetByTenantIDAndBranchID(c *gin.Context) {
	ctx, cancel := h.newTimeoutContext(c)
	defer cancel()

	tenantID, branchID, ok := h.getTenantAndBranchID(c)
	if !ok {
		return
	}

	data, err := h.posService.GetBranchDetail(ctx, tenantID, branchID)
	if err != nil {
		h.respondError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, data)
}

func (h *PosHandler) GetAllProducts(c *gin.Context) {
	ctx, cancel := h.newTimeoutContext(c)
	defer cancel()

	tenantID, branchID, ok := h.getTenantAndBranchID(c)
	if !ok {
		return
	}

	data, err := h.posService.GetProducts(ctx, tenantID, branchID)
	if err != nil {
		h.respondError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, data)
}

func (h *PosHandler) GetProductByID(c *gin.Context) {
	ctx, cancel := h.newTimeoutContext(c)
	defer cancel()

	tenantID, branchID, productID, ok := h.getTenantBranchAndProductID(c)
	if !ok {
		return
	}

	data, err := h.posService.GetProductByID(ctx, tenantID, branchID, productID)
	if err != nil {
		h.respondError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, toProductResponse(data))
}

func (h *PosHandler) CreateProduct(c *gin.Context) {
	ctx, cancel := h.newTimeoutContext(c)
	defer cancel()

	tenantID, branchID, ok := h.getTenantAndBranchID(c)
	if !ok {
		return
	}

	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.respondValidationError(c, err.Error())
		return
	}

	data, err := h.posService.CreateNewProduct(ctx, tenantID, branchID, req)
	if err != nil {
		h.respondError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, toProductResponse(data))
}

func (h *PosHandler) UpdateProduct(c *gin.Context) {
	ctx, cancel := h.newTimeoutContext(c)
	defer cancel()

	tenantID, branchID, productID, ok := h.getTenantBranchAndProductID(c)
	if !ok {
		return
	}

	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.respondValidationError(c, err.Error())
		return
	}

	data, err := h.posService.UpdateProduct(ctx, tenantID, branchID, productID, req)
	if err != nil {
		h.respondError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, toProductResponse(data))
}

func (h *PosHandler) DeleteProduct(c *gin.Context) {
	ctx, cancel := h.newTimeoutContext(c)
	defer cancel()

	tenantID, branchID, productID, ok := h.getTenantBranchAndProductID(c)
	if !ok {
		return
	}

	if err := h.posService.DeleteProduct(ctx, tenantID, branchID, productID); err != nil {
		h.respondError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "product deleted",
		"product_id": productID,
	})
}

func (h *PosHandler) getTenantAndBranchID(c *gin.Context) (string, string, bool) {
	tenantID, ok := h.getValidTenantID(c)
	if !ok {
		return "", "", false
	}

	branchID, ok := h.getValidBranchID(c)
	if !ok {
		return "", "", false
	}

	return tenantID, branchID, true
}

func (h *PosHandler) getTenantBranchAndProductID(c *gin.Context) (string, string, string, bool) {
	tenantID, branchID, ok := h.getTenantAndBranchID(c)
	if !ok {
		return "", "", "", false
	}

	productID, ok := h.getValidProductID(c)
	if !ok {
		return "", "", "", false
	}

	return tenantID, branchID, productID, true
}

func (h *PosHandler) getValidTenantID(c *gin.Context) (string, bool) {
	tenantID := c.Param("tenant_id")
	if tenantID == "" {
		h.respondValidationError(c, "tenant_id is required")
		return "", false
	}

	if err := h.validator.TenantIDValidation(tenantID); err != nil {
		h.respondValidationError(c, err.Error())
		return "", false
	}

	return tenantID, true
}

func (h *PosHandler) getValidBranchID(c *gin.Context) (string, bool) {
	branchID := c.Param("branch_id")
	if branchID == "" {
		h.respondValidationError(c, "branch_id is required")
		return "", false
	}

	if err := h.validator.BranchIDValidation(branchID); err != nil {
		h.respondValidationError(c, err.Error())
		return "", false
	}

	return branchID, true
}

func (h *PosHandler) getValidProductID(c *gin.Context) (string, bool) {
	productID := c.Param("product_id")
	if productID == "" {
		h.respondValidationError(c, "product_id is required")
		return "", false
	}

	if err := h.validator.ProductIDValidation(productID); err != nil {
		h.respondValidationError(c, err.Error())
		return "", false
	}

	return productID, true
}

func (h *PosHandler) Readiness(c *gin.Context) {
	if err := h.posService.GetHealth(c.Request.Context()); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "not_ready",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
	})
}