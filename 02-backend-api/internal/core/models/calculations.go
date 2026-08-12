package models

import "time"

type Calculation struct {
	ID        int64
	Input     []float64
	Operation string
	Result    float64
	Username  string
	CreatedAt time.Time
}
