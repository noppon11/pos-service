package domain

import "time"

type Product struct {
	ProductID  string
	Name       string
	SKU        string
	Price      float64
	CategoryID string
	Unit       string
	IsActive   bool
	DeletedAt  *time.Time
}