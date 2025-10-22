package store

import "database/sql"

type Store struct {
	db *sql.DB
}

// create a new store by injecting db
func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}
