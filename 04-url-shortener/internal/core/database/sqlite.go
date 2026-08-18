package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

func Open(dbPath string) (*sql.DB, error) {
	// Create Directory if not exists
	if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
		return nil, fmt.Errorf("creating db dir: %w", err)
	}

	dsn := fmt.Sprintf("%s?_journal_mode=WAL&_busy_timeout=5000&_foreign_keys=ON", dbPath)
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Ping verifies the connection is actually valid
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Create tables if they don't exist
	schema := `
	CREATE TABLE urls (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		short_code TEXT NOT NULL UNIQUE COLLATE BINARY,
		original_url TEXT NOT NULL,
		clicks INTEGER NOT NULL DEFAULT 0,
		user_id TEXT DEFAULT NULL,
		is_active INTEGER NOT NULL DEFAULT 1,
		expires_at DATETIME DEFAULT NULL,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX idx_urls_short_code ON urls(short_code);
	CREATE INDEX idx_urls_user_id ON urls(user_id);`

	if _, err := db.Exec(schema); err != nil {
		return nil, fmt.Errorf("failed to execute schema: %w", err)
	}

	return db, nil
}
