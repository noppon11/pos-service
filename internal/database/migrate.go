package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func Migrate(db *sql.DB) {
	if err := ensureSchemaMigrationsTable(db); err != nil {
		log.Fatalf("failed to ensure schema_migrations table: %v", err)
	}

	migrationDir := "internal/database/migrations"

	files, err := os.ReadDir(migrationDir)
	if err != nil {
		log.Fatalf("failed to read migration dir: %v", err)
	}

	var upFiles []string
	for _, f := range files {
		name := f.Name()
		if strings.HasSuffix(name, ".up.sql") {
			upFiles = append(upFiles, name)
		}
	}

	sort.Strings(upFiles)

	for _, file := range upFiles {
		version := strings.Split(file, "_")[0]

		applied, err := isMigrationApplied(db, version)
		if err != nil {
			log.Fatalf("failed checking migration version %s: %v", version, err)
		}
		if applied {
			continue
		}

		path := filepath.Join(migrationDir, file)
		content, err := os.ReadFile(path)
		if err != nil {
			log.Fatalf("failed to read migration file %s: %v", file, err)
		}

		tx, err := db.Begin()
		if err != nil {
			log.Fatalf("failed to begin tx for %s: %v", file, err)
		}

		if _, err := tx.Exec(string(content)); err != nil {
			_ = tx.Rollback()
			log.Fatalf("failed to execute migration %s: %v", file, err)
		}

		if _, err := tx.Exec(
			`INSERT INTO schema_migrations (version) VALUES ($1)`,
			version,
		); err != nil {
			_ = tx.Rollback()
			log.Fatalf("failed to record migration %s: %v", file, err)
		}

		if err := tx.Commit(); err != nil {
			log.Fatalf("failed to commit migration %s: %v", file, err)
		}

		log.Printf("applied migration: %s", file)
	}

	log.Println("all migrations applied")
}

func ensureSchemaMigrationsTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS schema_migrations (
		version TEXT PRIMARY KEY,
		applied_at TIMESTAMP NOT NULL DEFAULT NOW()
	);
	`
	_, err := db.Exec(query)
	return err
}

func isMigrationApplied(db *sql.DB, version string) (bool, error) {
	var exists bool
	err := db.QueryRow(
		`SELECT EXISTS (SELECT 1 FROM schema_migrations WHERE version = $1)`,
		version,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("query migration version: %w", err)
	}
	return exists, nil
}