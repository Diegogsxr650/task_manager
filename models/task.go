package models

import "time"

type Task struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at"` // Puntero para permitir NULL
}

func (t *Task) MarkAsCompleted() {
	t.Status = "completed"
	now := time.Now()
	t.CompletedAt = &now // Assign the current timestamp
}
