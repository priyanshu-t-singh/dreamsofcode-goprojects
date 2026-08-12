package arithmetics

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/priyanshu-t-singh/dreamsofcode-goprojects/02-backend-api/internal/core/models"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) SaveCalculation(calc *models.Calculation) error {
	query := `
	INSERT INTO calculations (input, operation, result, username, created_at)
	VALUES (?, ?, ?, ?, ?)`

	input, err := json.Marshal(calc.Input)
	if err != nil {
		return fmt.Errorf("failed to parse the input array")
	}

	_, err2 := r.db.Exec(query, string(input), calc.Operation, calc.Result, calc.Username, time.Now().UTC())
	if err2 != nil {
		return fmt.Errorf("failed to save calculation: %w", err)
	}

	return nil
}
