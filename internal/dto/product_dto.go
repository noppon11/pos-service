package dto

import (
	"math"
	"pos-service/internal/domain"
	"time"
)

type CreateProductRequest struct {
	Name       string  `json:"name" binding:"required"`
	SKU        string  `json:"sku" binding:"required"`
	Price      float64 `json:"price" binding:"required"`
	CategoryID string  `json:"category_id" binding:"required"`
	Unit       string  `json:"unit" binding:"required"`
	IsActive   bool    `json:"is_active"`
}

type UpdateProductRequest struct {
	Name       string  `json:"name" binding:"required"`
	SKU        string  `json:"sku" binding:"required"`
	Price      float64 `json:"price" binding:"required"`
	CategoryID string  `json:"category_id" binding:"required"`
	Unit       string  `json:"unit" binding:"required"`
	IsActive   bool    `json:"is_active"`
}

type ListProductsQuery struct {
	Page       int    `form:"page"`
	Limit      int    `form:"limit"`
	CategoryID string `form:"category_id"`
}

type ProductResponse struct {
	ProductID  string  `json:"product_id"`
	Name       string  `json:"name"`
	SKU        string  `json:"sku"`
	Price      float64 `json:"price"`
	CategoryID string  `json:"category_id"`
	Unit       string  `json:"unit"`
	IsActive   bool    `json:"is_active"`
	DeletedAt  *string `json:"deleted_at,omitempty"`
}

type PaginationMeta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}

type ListProductsResponse struct {
	Items []ProductResponse `json:"items"`
	Meta  PaginationMeta    `json:"meta"`
}

func ToProductResponse(p domain.Product) ProductResponse {
	var deletedAt *string
	if p.DeletedAt != nil {
		v := p.DeletedAt.Format(time.RFC3339)
		deletedAt = &v
	}

	return ProductResponse{
		ProductID:  p.ProductID,
		Name:       p.Name,
		SKU:        p.SKU,
		Price:      p.Price,
		CategoryID: p.CategoryID,
		Unit:       p.Unit,
		IsActive:   p.IsActive,
		DeletedAt:  deletedAt,
	}
}

func ToProductResponsePtr(p *domain.Product) *ProductResponse {
	if p == nil {
		return nil
	}
	resp := ToProductResponse(*p)
	return &resp
}

func ToListProductsResponse(products []domain.Product, page, limit, total int) ListProductsResponse {
	items := make([]ProductResponse, 0, len(products))
	for _, p := range products {
		items = append(items, ToProductResponse(p))
	}

	totalPages := 0
	if limit > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(limit)))
	}

	return ListProductsResponse{
		Items: items,
		Meta: PaginationMeta{
			Page:       page,
			Limit:      limit,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}
}