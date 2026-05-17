package store

import "github.com/jmoiron/sqlx"

// Store is the single access point for all database operations.
// Each domain (courses, topics, notes, etc.) gets its own *_store.go
// file with methods on this struct.
type Store struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Store {
	return &Store{db: db}
}
