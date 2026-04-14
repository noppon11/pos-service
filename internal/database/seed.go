package database

import (
	"database/sql"
	"log"
)

func Seed(db *sql.DB) {
	query := `
	INSERT INTO branches (tenant_id, branch_id, branch_name, status, timezone, currency)
	VALUES
	('aura-bkk', 'bkk-001', 'Aura Siam', 'active', 'Asia/Bangkok', 'THB'),
	('aura-bkk', 'bkk-002', 'Aura Ari', 'inactive', 'Asia/Bangkok', 'THB')
	ON CONFLICT (tenant_id, branch_id) DO NOTHING;
	`

	if _, err := db.Exec(query); err != nil {
		log.Fatalf("failed to seed: %v", err)
	}

	log.Println("seed done")
}