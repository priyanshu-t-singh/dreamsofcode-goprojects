package tasks

import (
	"database/sql"
	"fmt"
	"time"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(title string, dueDate *string) (int64, error) {
	query := `INSERT INTO tasks (title, created_at, due_date) VALUES (?, ?, ?);`
	createdAt := time.Now().UTC()
	result, err := r.db.Exec(query, title, createdAt, dueDate)
	if err != nil {
		return 0, fmt.Errorf("failed to create task: %w", err)
	}

	return result.LastInsertId()
}

func (r *Repository) GetAll(showAll bool) ([]Task, error) {
	query := "SELECT id, title, created_at, completed_at, due_date FROM tasks\n"
	if !showAll {
		query += "WHERE completed_at IS NULL"
	}

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Title, &task.CreatedAt, &task.CompletedAt, &task.DueDate)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *Repository) MarkAsCompleted(taskId string) error {
	query := `UPDATE tasks SET completed_at = ? WHERE id = ?`
	completedAt := time.Now().UTC()
	_, err := r.db.Exec(query, completedAt, taskId)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	return nil
}

// func (r *Repository) Delete() (int, error) {

// }
