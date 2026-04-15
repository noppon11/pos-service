package dto

type CreateProductRequest struct {
	Name       string  `json:"name" binding:"required"`
	SKU        string  `json:"sku" binding:"required"`
	Price      float64 `json:"price" binding:"required"`
	CategoryID string  `json:"category_id" binding:"required"`
	Unit       string  `json:"unit" binding:"required"`
	IsActive   bool    `json:"is_active"`
}

type UpdateProductRequest struct {
	ProductID   string  `json:"product_id" binding:"required"`
	Name       string  `json:"name" binding:"required"`
	SKU        string  `json:"sku" binding:"required"`
	Price      float64 `json:"price" binding:"required"`
	CategoryID string  `json:"category_id" binding:"required"`
	Unit       string  `json:"unit" binding:"required"`
	IsActive   bool    `json:"is_active"`
}