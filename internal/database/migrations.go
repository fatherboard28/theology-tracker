package database

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"sort"
	"strings"

	"github.com/jmoiron/sqlx"
)

//go:embed migrations
var migrationsFS embed.FS

func RunMigrations(db *sqlx.DB) error {
	// The migrations table itself is created outside the migration system
	// so the system can always bootstrap regardless of state.
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version    TEXT PRIMARY KEY,
			applied_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("creating schema_migrations table: %w", err)
	}

	entries, err := fs.ReadDir(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("reading migrations directory: %w", err)
	}

	var files []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".sql") {
			files = append(files, e.Name())
		}
	}
	sort.Strings(files)

	for _, file := range files {
		version := strings.TrimSuffix(file, ".sql")

		var count int
		if err := db.QueryRow(
			"SELECT COUNT(*) FROM schema_migrations WHERE version = ?", version,
		).Scan(&count); err != nil {
			return fmt.Errorf("checking migration %s: %w", version, err)
		}
		if count > 0 {
			continue
		}

		content, err := migrationsFS.ReadFile("migrations/" + file)
		if err != nil {
			return fmt.Errorf("reading migration %s: %w", file, err)
		}

		if _, err := db.Exec(string(content)); err != nil {
			return fmt.Errorf("applying migration %s: %w", file, err)
		}

		if _, err := db.Exec(
			"INSERT INTO schema_migrations (version) VALUES (?)", version,
		); err != nil {
			return fmt.Errorf("recording migration %s: %w", file, err)
		}

		log.Printf("migration applied: %s", version)
	}

	return nil
}
