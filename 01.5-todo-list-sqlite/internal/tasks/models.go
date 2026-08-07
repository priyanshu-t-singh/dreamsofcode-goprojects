package tasks

import "time"

type Task struct {
	ID          string
	Title       string
	CreatedAt   time.Time
	CompletedAt *time.Time
	DueDate     *time.Time
}
