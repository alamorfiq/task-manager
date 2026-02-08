package dto

import "taskmgr/internal/storage"

// CreateTaskRequest - запрос на создание задачи
type CreateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// Validate проверяет валидность запроса
func (r *CreateTaskRequest) Validate() error {
	if r.Title == "" {
		return ErrTitleRequired
	}
	return nil
}

// TaskResponse - ответ с одной задачей
type TaskResponse struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	IsDone      bool   `json:"is_done"`
	CreatedAt   string `json:"created_at"`
}

// TasksResponse - ответ со списком задач
type TasksResponse struct {
	Tasks []TaskResponse `json:"tasks"`
	Count int            `json:"count"`
}

// ToTaskResponse преобразует storage.Task в TaskResponse
func ToTaskResponse(t storage.Task) TaskResponse {
	return TaskResponse{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		IsDone:      t.IsDone,
		CreatedAt:   t.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// ToTasksResponse преобразует список задач в TasksResponse
func ToTasksResponse(tasks []storage.Task) TasksResponse {
	responses := make([]TaskResponse, len(tasks))
	for i, task := range tasks {
		responses[i] = ToTaskResponse(task)
	}
	return TasksResponse{
		Tasks: responses,
		Count: len(responses),
	}
}

