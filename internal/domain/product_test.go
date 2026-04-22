package domain

import (
	"testing"
	"time"
)

func TestProduct_IsAvailable(t *testing.T) {
	tests := []struct {
		name     string
		product  Product
		expected bool
	}{
		{
			name: "สินค้าพร้อมขาย — active และมีสต็อก",
			product: Product{
				ProductID: "prod-001",
				Name:      "Facial Cream",
				SKU:       "FC-001",
				Price:     59900,
				Stock:     10,
				IsActive:  true,
				DeletedAt: nil,
			},
			expected: true,
		},
		{
			name: "สินค้าไม่ active — ไม่ขาย",
			product: Product{
				ProductID: "prod-002",
				Name:      "Old Product",
				SKU:       "OLD-001",
				Stock:     5,
				IsActive:  false,
				DeletedAt: nil,
			},
			expected: false,
		},
		{
			name: "สินค้าหมดสต็อก — ไม่ขาย",
			product: Product{
				ProductID: "prod-003",
				Name:      "Out of Stock",
				SKU:       "OOS-001",
				Stock:     0,
				IsActive:  true,
				DeletedAt: nil,
			},
			expected: false,
		},
		{
			name: "สินค้าถูก soft delete — ไม่ขาย",
			product: Product{
				ProductID: "prod-004",
				Name:      "Deleted Product",
				SKU:       "DEL-001",
				Stock:     10,
				IsActive:  true,
				DeletedAt: timePtr(time.Now()),
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.product.IsAvailable()
			if got != tt.expected {
				t.Errorf("IsAvailable() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestProduct_IsLowStock(t *testing.T) {
	tests := []struct {
		name     string
		stock    int64
		expected bool
	}{
		{"สต็อกต่ำ — 5 ชิ้น", 5, true},
		{"สต็อกต่ำ — 9 ชิ้น", 9, true},
		{"สต็อกพอดี — 10 ชิ้น", 10, false},
		{"สต็อกเยอะ — 50 ชิ้น", 50, false},
		{"สต็อกหมด — 0 ชิ้น", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			product := Product{Stock: tt.stock}
			got := product.IsLowStock()
			if got != tt.expected {
				t.Errorf("IsLowStock() with stock=%d = %v, want %v", tt.stock, got, tt.expected)
			}
		})
	}
}

func TestProduct_DeductStock(t *testing.T) {
	tests := []struct {
		name          string
		initialStock  int64
		deductQty     int64
		expectedStock int64
		expectError   bool
		errorMessage  string
	}{
		{
			name:          "หักสต็อกปกติ — 10 - 3 = 7",
			initialStock:  10,
			deductQty:     3,
			expectedStock: 7,
			expectError:   false,
		},
		{
			name:          "หักสต็อกพอดี — 5 - 5 = 0",
			initialStock:  5,
			deductQty:     5,
			expectedStock: 0,
			expectError:   false,
		},
		{
			name:          "หักสต็อกมากกว่าที่มี — error",
			initialStock:  3,
			deductQty:     5,
			expectedStock: 3, // ไม่เปลี่ยน
			expectError:   true,
			errorMessage:  "insufficient stock",
		},
		{
			name:          "หักสต็อกด้วยจำนวนติดลบ — error",
			initialStock:  10,
			deductQty:     -2,
			expectedStock: 10, // ไม่เปลี่ยน
			expectError:   true,
			errorMessage:  "quantity must be greater than 0",
		},
		{
			name:          "หักสต็อกด้วย 0 — error",
			initialStock:  10,
			deductQty:     0,
			expectedStock: 10, // ไม่เปลี่ยน
			expectError:   true,
			errorMessage:  "quantity must be greater than 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			product := &Product{
				ProductID: "test-prod",
				Stock:     tt.initialStock,
			}

			err := product.DeductStock(tt.deductQty)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got nil")
					return
				}
				// ตรวจสอบ error message (case-insensitive)
				if tt.errorMessage != "" && !contains(err.Error(), tt.errorMessage) {
					t.Errorf("expected error containing '%s', got '%s'", tt.errorMessage, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					return
				}
			}

			if product.Stock != tt.expectedStock {
				t.Errorf("Stock = %d, want %d", product.Stock, tt.expectedStock)
			}
		})
	}
}

func TestProduct_UpdatePrice(t *testing.T) {
	tests := []struct {
		name        string
		oldPrice    int64
		newPrice    int64
		expectError bool
	}{
		{
			name:        "เปลี่ยนราคาปกติ — 599 → 699",
			oldPrice:    59900,
			newPrice:    69900,
			expectError: false,
		},
		{
			name:        "ลดราคา — 1000 → 799",
			oldPrice:    100000,
			newPrice:    79900,
			expectError: false,
		},
		{
			name:        "ราคาเป็น 0 — error",
			oldPrice:    59900,
			newPrice:    0,
			expectError: true,
		},
		{
			name:        "ราคาติดลบ — error",
			oldPrice:    59900,
			newPrice:    -100,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			product := &Product{
				ProductID: "test-prod",
				Price:     tt.oldPrice,
			}

			err := product.UpdatePrice(tt.newPrice)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got nil")
				}
				if product.Price != tt.oldPrice {
					t.Errorf("Price changed when error expected: %d → %d", tt.oldPrice, product.Price)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if product.Price != tt.newPrice {
					t.Errorf("Price = %d, want %d", product.Price, tt.newPrice)
				}
			}
		})
	}
}

func TestProduct_Activate(t *testing.T) {
	product := &Product{
		ProductID: "test-prod",
		IsActive:  false,
	}

	product.Activate()

	if !product.IsActive {
		t.Errorf("IsActive = false, want true")
	}
}

func TestProduct_Deactivate(t *testing.T) {
	product := &Product{
		ProductID: "test-prod",
		IsActive:  true,
	}

	product.Deactivate()

	if product.IsActive {
		t.Errorf("IsActive = true, want false")
	}
}

func TestProduct_SoftDelete(t *testing.T) {
	product := &Product{
		ProductID: "test-prod",
		DeletedAt: nil,
	}

	now := time.Now()
	product.SoftDelete(now)

	if product.DeletedAt == nil {
		t.Errorf("DeletedAt = nil, want %v", now)
	}
	if !product.DeletedAt.Equal(now) {
		t.Errorf("DeletedAt = %v, want %v", product.DeletedAt, now)
	}
}

func TestProduct_IsDeleted(t *testing.T) {
	tests := []struct {
		name      string
		deletedAt *time.Time
		expected  bool
	}{
		{
			name:      "ยังไม่ถูกลบ",
			deletedAt: nil,
			expected:  false,
		},
		{
			name:      "ถูกลบแล้ว",
			deletedAt: timePtr(time.Now()),
			expected:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			product := &Product{
				DeletedAt: tt.deletedAt,
			}
			got := product.IsDeleted()
			if got != tt.expected {
				t.Errorf("IsDeleted() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func timePtr(t time.Time) *time.Time {
	return &t
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || containsStr(s, substr))
}

func containsStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}