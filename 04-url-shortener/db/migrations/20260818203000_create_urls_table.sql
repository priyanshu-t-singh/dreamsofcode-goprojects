-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS urls (
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

CREATE INDEX IF NOT EXISTS idx_urls_short_code ON urls(short_code);
CREATE INDEX IF NOT EXISTS idx_urls_user_id ON urls(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_urls_user_id;
DROP INDEX IF EXISTS idx_urls_short_code;
DROP TABLE IF EXISTS urls;
-- +goose StatementEnd
