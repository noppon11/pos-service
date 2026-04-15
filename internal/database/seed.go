package database

import (
	"database/sql"
	"log"
)

func Seed(db *sql.DB) {
	queries := []string{
		`
		INSERT INTO branches (tenant_id, branch_id, branch_name, status, timezone, currency)
		VALUES
			('aura-bkk', 'bkk-001', 'Aura Siam', 'active', 'Asia/Bangkok', 'THB'),
			('aura-bkk', 'bkk-002', 'Aura Ari', 'inactive', 'Asia/Bangkok', 'THB')
		ON CONFLICT (tenant_id, branch_id) DO NOTHING;
		`,
		`
		INSERT INTO products (
			product_id, tenant_id, branch_id, name, sku, price, category_id, unit, is_active
		)
		VALUES
			('prod-001', 'aura-bkk', 'bkk-001', 'Botox 50u', 'BOT-50', 3500, 'treatment', 'unit', true),
			('prod-002', 'aura-bkk', 'bkk-001', 'Filler 1cc', 'FIL-1', 12000, 'treatment', 'unit', true),
			('prod-003', 'aura-bkk', 'bkk-002', 'Laser Acne', 'LAS-ACNE', 2500, 'treatment', 'session', false)
		ON CONFLICT (tenant_id, branch_id, product_id) DO NOTHING;
		`,
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("failed to begin seed tx: %v", err)
	}
	defer tx.Rollback()

	for _, q := range queries {
		if _, err := tx.Exec(q); err != nil {
			log.Fatalf("failed to seed: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		log.Fatalf("failed to commit seed: %v", err)
	}

	log.Println("seed done")
}