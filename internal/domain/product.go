package domain

import (
	"time"
	appErr "pos-service/internal/errors"
)

type Product struct {
	ProductID  string
	Name       string
	SKU        string
	Price      int64
	CategoryID string
	Unit       string
	Stock      int64
	Version    int64
	IsActive   bool
	DeletedAt  *time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (p *Product) IsAvailable() bool {
	return p.IsActive && p.Stock > 0 && p.DeletedAt == nil
}

func (p *Product) IsLowStock() bool {
	return p.Stock < 10
}

func (p *Product) DeductStock(qty int64) error {
	if qty <= 0 {
		return appErr.ErrInvalidQuantity
	}
	if p.Stock < qty {
		return appErr.ErrInsufficientStock
	}
	p.Stock -= qty
	p.UpdatedAt = time.Now()
	return nil
}

func (p *Product) AddStock(qty int64) error {
	if qty <= 0 {
		return appErr.ErrInvalidQuantity
	}
	p.Stock += qty
	p.UpdatedAt = time.Now()
	return nil
}

func (p *Product) UpdatePrice(newPrice int64) error {
	if newPrice <= 0 {
		return appErr.ErrInvalidPrice
	}
	p.Price = newPrice
	p.UpdatedAt = time.Now()
	return nil
}

func (p *Product) Activate() {
	p.IsActive = true
	p.UpdatedAt = time.Now()
}

func (p *Product) Deactivate() {
	p.IsActive = false
	p.UpdatedAt = time.Now()
}

func (p *Product) SoftDelete(at time.Time) {
	p.DeletedAt = &at
	p.IsActive = false
	p.UpdatedAt = time.Now()
}

func (p *Product) IsDeleted() bool {
	return p.DeletedAt != nil
}