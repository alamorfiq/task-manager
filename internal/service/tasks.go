package service

import (
	"context"
	"taskmgr/internal/storage"
)

type TaskService struct {
	storage storage.Storage
}

func NewTaskService(s storage.Storage) *TaskService {
	return &TaskService{storage: s}
}

func (s *TaskService) Create(ctx context.Context, title, desc string) (int64, error) {
	return s.storage.CreateTask(ctx, storage.Task{
		Title:       title,
		Description: desc,
	})
}

func (s *TaskService) Get(ctx context.Context, id int64) (storage.Task, error) {
	return s.storage.GetTask(ctx, id)
}

func (s *TaskService) List(ctx context.Context) ([]storage.Task, error) {
	return s.storage.ListTasks(ctx)
}

func (s *TaskService) MarkDone(ctx context.Context, id int64) error {
	return s.storage.MarkDone(ctx, id)
}

func (s *TaskService) Delete(ctx context.Context, id int64) error {
	return s.storage.DeleteTask(ctx, id)
}
