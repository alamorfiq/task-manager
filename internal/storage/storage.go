package storage

import "context"

// Storage — базовый интерфейс
type Storage interface {
	Close() error
	Ping(ctx context.Context) error

	CreateTask(ctx context.Context, t Task) (int64, error)
	GetTask(ctx context.Context, id int64) (Task, error)
	ListTasks(ctx context.Context) ([]Task, error)
	MarkDone(ctx context.Context, id int64) error
	DeleteTask(ctx context.Context, id int64) error
}
