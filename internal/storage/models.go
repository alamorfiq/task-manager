package storage

import "time"

type Task struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	IsDone      bool      `json:"is_done"`
	CreatedAt   time.Time `json:"created_at"`
}
