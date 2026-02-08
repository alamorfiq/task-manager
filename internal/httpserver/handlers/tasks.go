package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"

	"taskmgr/internal/httpserver/dto"
	"taskmgr/internal/httpserver/response"
	"taskmgr/internal/service"
)

type TaskHandler struct {
	service *service.TaskService
	log     *slog.Logger
}

func NewTaskHandler(service *service.TaskService, log *slog.Logger) *TaskHandler {
	return &TaskHandler{
		service: service,
		log:     log,
	}
}

// CreateTask создает новую задачу
// POST /api/v1/tasks
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.TaskHandler.CreateTask"
	log := h.log.With(slog.String("op", op))

	var req dto.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("failed to decode request", slog.String("error", err.Error()))
		response.BadRequest(w, "invalid request body")
		return
	}

	if err := req.Validate(); err != nil {
		log.Warn("validation failed", slog.String("error", err.Error()))
		response.BadRequest(w, err.Error())
		return
	}

	id, err := h.service.Create(r.Context(), req.Title, req.Description)
	if err != nil {
		log.Error("failed to create task", slog.String("error", err.Error()))
		response.InternalError(w, "failed to create task")
		return
	}

	task, err := h.service.Get(r.Context(), id)
	if err != nil {
		log.Error("failed to get created task", slog.String("error", err.Error()))
		response.InternalError(w, "task created but failed to fetch")
		return
	}

	log.Info("task created", slog.Int64("id", id))
	response.Created(w, dto.ToTaskResponse(task))
}

// GetTask возвращает задачу по ID
// GET /api/v1/tasks/{id}
func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.TaskHandler.GetTask"
	log := h.log.With(slog.String("op", op))

	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Warn("invalid task id", slog.String("id", idParam))
		response.BadRequest(w, "invalid task id")
		return
	}

	task, err := h.service.Get(r.Context(), id)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Info("task not found", slog.Int64("id", id))
			response.NotFound(w, "task not found")
			return
		}
		log.Error("failed to get task", slog.String("error", err.Error()))
		response.InternalError(w, "failed to get task")
		return
	}

	response.OK(w, dto.ToTaskResponse(task))
}

// ListTasks возвращает список всех задач
// GET /api/v1/tasks
func (h *TaskHandler) ListTasks(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.TaskHandler.ListTasks"
	log := h.log.With(slog.String("op", op))

	tasks, err := h.service.List(r.Context())
	if err != nil {
		log.Error("failed to list tasks", slog.String("error", err.Error()))
		response.InternalError(w, "failed to list tasks")
		return
	}

	response.OK(w, dto.ToTasksResponse(tasks))
}

// MarkTaskDone отмечает задачу как выполненную
// PATCH /api/v1/tasks/{id}/done
func (h *TaskHandler) MarkTaskDone(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.TaskHandler.MarkTaskDone"
	log := h.log.With(slog.String("op", op))

	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Warn("invalid task id", slog.String("id", idParam))
		response.BadRequest(w, "invalid task id")
		return
	}

	// Проверяем что задача существует
	_, err = h.service.Get(r.Context(), id)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Info("task not found", slog.Int64("id", id))
			response.NotFound(w, "task not found")
			return
		}
		log.Error("failed to get task", slog.String("error", err.Error()))
		response.InternalError(w, "failed to get task")
		return
	}

	if err := h.service.MarkDone(r.Context(), id); err != nil {
		log.Error("failed to mark task as done", slog.String("error", err.Error()))
		response.InternalError(w, "failed to mark task as done")
		return
	}

	// Получаем обновленную задачу
	task, err := h.service.Get(r.Context(), id)
	if err != nil {
		log.Error("failed to get updated task", slog.String("error", err.Error()))
		response.InternalError(w, "task updated but failed to fetch")
		return
	}

	log.Info("task marked as done", slog.Int64("id", id))
	response.OK(w, dto.ToTaskResponse(task))
}

// DeleteTask удаляет задачу
// DELETE /api/v1/tasks/{id}
func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.TaskHandler.DeleteTask"
	log := h.log.With(slog.String("op", op))

	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Warn("invalid task id", slog.String("id", idParam))
		response.BadRequest(w, "invalid task id")
		return
	}

	// Проверяем что задача существует
	_, err = h.service.Get(r.Context(), id)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Info("task not found", slog.Int64("id", id))
			response.NotFound(w, "task not found")
			return
		}
		log.Error("failed to get task", slog.String("error", err.Error()))
		response.InternalError(w, "failed to get task")
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		log.Error("failed to delete task", slog.String("error", err.Error()))
		response.InternalError(w, "failed to delete task")
		return
	}

	log.Info("task deleted", slog.Int64("id", id))
	response.OK(w, map[string]string{"message": "task deleted successfully"})
}

