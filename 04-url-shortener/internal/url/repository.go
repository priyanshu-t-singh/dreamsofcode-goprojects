package url

import (
	"database/sql"
	"errors"
	"fmt"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) SaveShortURL(c *CreateShortURLSQL) error {
	query := `
	INSERT INTO urls (short_code, original_url)
	VALUES (?, ?);`

	_, err := r.db.Exec(query, c.ShortCode, c.OriginalURL)
	if err != nil {
		return fmt.Errorf("failed to save short url: %w", err)
	}

	return nil
}

func (r *Repository) GetOriginalURL(id string) (string, error) {
	query := `SELECT original_url FROM urls WHERE short_code = ?`

	var u string
	resp := r.db.QueryRow(query, id)
	if err := resp.Scan(&u); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "-1", fmt.Errorf("url not found")
		}
		return "", err
	}

	return u, nil
}
