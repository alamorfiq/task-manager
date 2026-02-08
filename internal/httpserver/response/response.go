package response

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// Response - стандартный формат ответа
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// JSON отправляет JSON ответ
func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp := Response{
		Success: statusCode < 400,
		Data:    data,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error("failed to encode response", slog.String("error", err.Error()))
	}
}

// Error отправляет ответ с ошибкой
func Error(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp := Response{
		Success: false,
		Error:   message,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error("failed to encode error response", slog.String("error", err.Error()))
	}
}

// OK отправляет успешный ответ с данными
func OK(w http.ResponseWriter, data interface{}) {
	JSON(w, http.StatusOK, data)
}

// Created отправляет ответ при создании ресурса
func Created(w http.ResponseWriter, data interface{}) {
	JSON(w, http.StatusCreated, data)
}

// BadRequest отправляет ответ с ошибкой 400
func BadRequest(w http.ResponseWriter, message string) {
	Error(w, http.StatusBadRequest, message)
}

// NotFound отправляет ответ с ошибкой 404
func NotFound(w http.ResponseWriter, message string) {
	Error(w, http.StatusNotFound, message)
}

// InternalError отправляет ответ с ошибкой 500
func InternalError(w http.ResponseWriter, message string) {
	Error(w, http.StatusInternalServerError, message)
}

