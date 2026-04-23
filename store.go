// Package hookstore is the registry for bridge-managed harness hooks.
//
// Each row is a hook handler (harness, event, matcher, shell command) bound
// to a scope (global, instance, session). The llm-bridge server queries this
// store at session spawn to decide which hooks to wire into the harness's
// native hook mechanism (for Claude Code: synthesized into a settings JSON
// passed via --settings).
//
// The store itself is a pure data layer — it does not execute hooks, render
// settings JSON, or talk to harnesses. That belongs to the server.
package hookstore

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

// Store wraps the hook-store SQLite database.
type Store struct {
	db *sql.DB
}

// Open opens or creates a hook-store database at dbPath. The parent directory
// is created if missing. WAL is enabled for concurrent reads during writes.
func Open(dbPath string) (*Store, error) {
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("create db dir: %w", err)
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec("PRAGMA journal_mode=WAL; PRAGMA busy_timeout=5000; PRAGMA foreign_keys=ON;"); err != nil {
		db.Close()
		return nil, fmt.Errorf("sqlite pragmas: %w", err)
	}

	s := &Store{db: db}
	if err := s.migrate(); err != nil {
		db.Close()
		return nil, err
	}

	return s, nil
}

// Close closes the database connection.
func (s *Store) Close() error {
	return s.db.Close()
}

// DB returns the underlying database connection for advanced queries.
func (s *Store) DB() *sql.DB {
	return s.db
}

func (s *Store) migrate() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS hooks (
			id          TEXT PRIMARY KEY,
			harness     TEXT NOT NULL,
			event       TEXT NOT NULL,
			matcher     TEXT NOT NULL DEFAULT '',
			command     TEXT NOT NULL,
			scope_kind  TEXT NOT NULL,
			scope_id    TEXT NOT NULL DEFAULT '',
			enabled     INTEGER NOT NULL DEFAULT 1,
			created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			CHECK (scope_kind IN ('global', 'instance', 'session')),
			CHECK (scope_kind = 'global' OR scope_id != '')
		);
		CREATE INDEX IF NOT EXISTS idx_hooks_harness_event ON hooks(harness, event);
		CREATE INDEX IF NOT EXISTS idx_hooks_scope ON hooks(scope_kind, scope_id);
		CREATE INDEX IF NOT EXISTS idx_hooks_enabled ON hooks(enabled);
	`)
	return err
}
