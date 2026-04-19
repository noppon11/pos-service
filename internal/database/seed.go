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
		`
		INSERT INTO users (
			id, email, password_hash, full_name, tenant_id, role, is_active
		)
		VALUES (
			'usr_001',
			'staff@aura.com',
			'$2a$10$wQmWm3P7QJ3k7KXq5v6jAeg7lP4b4vA0kYYmK6k0iK0c6xWQnQG2K',
			'Aura Staff',
			'aura-bkk',
			'staff',
			TRUE
		)
		ON CONFLICT (id) DO NOTHING;
		`,
		`
		INSERT INTO user_branch_access (user_id, branch_id)
		VALUES ('usr_001', 'bkk-001')
		ON CONFLICT (user_id, branch_id) DO NOTHING;
		`,
		`
		INSERT INTO users (
			id, email, password_hash, full_name, tenant_id, role, is_active
		)
		VALUES (
			'usr_001',
			'staff@aura.com',
			'$2a$10$.N9oLtfq49Jct2IZanv7n.9r7jL8GvyvQ9Pf4g59VFf/SbYTH1E/K',
			'Aura Staff',
			'aura-bkk',
			'staff',
			TRUE
		)
		ON CONFLICT (id) DO UPDATE SET
			password_hash = EXCLUDED.password_hash;
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
