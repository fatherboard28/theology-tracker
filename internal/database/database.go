package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

func Open(path string) (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}

	// SQLite performs best with a single writer connection.
	// Concurrent reads still work via WAL mode below.
	db.SetMaxOpenConns(1)

	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		return nil, fmt.Errorf("enabling WAL mode: %w", err)
	}

	if _, err := db.Exec("PRAGMA foreign_keys=ON"); err != nil {
		return nil, fmt.Errorf("enabling foreign keys: %w", err)
	}

	// Reduce the chance of SQLITE_BUSY under light concurrent load.
	if _, err := db.Exec("PRAGMA busy_timeout=5000"); err != nil {
		return nil, fmt.Errorf("setting busy timeout: %w", err)
	}

	return db, nil
}
