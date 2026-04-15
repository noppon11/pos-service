package domain

import "time"

type ProductResponse struct {
	ProductID   string  `json:"product_id"`
	Name string `json:"name"`
	SKU  string `json:"sku"`
	Price float64 `json:"price"`
	CategoryID string `json:"category_id"`
	Unit string `json:"unit"`
	IsActive bool `json:"is_active"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
}