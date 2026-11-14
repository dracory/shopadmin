package main

import "testing"

// openDB returns a usable in-memory SQLite connection.
func TestOpenDB(t *testing.T) {
	db, err := openDB(":memory:")
	if err != nil {
		t.Fatalf("openDB: %v", err)
	}
	defer func() { _ = db.Close() }()

	if err := db.Ping(); err != nil {
		t.Fatalf("ping: %v", err)
	}
}
