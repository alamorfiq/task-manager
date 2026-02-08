package httpserver

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"taskmgr/internal/httpserver/handlers"
	"taskmgr/internal/httpserver/response"
	"taskmgr/internal/service"
)

// NewRouter создает и настраивает роутер
func NewRouter(taskService *service.TaskService, log *slog.Logger) http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(LoggerMiddleware(log))
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// CORS для фронтенда
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://127.0.0.1:3000", "http://localhost:8082"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Статические файлы (фронтенд)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/index.html")
	})

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		response.OK(w, map[string]string{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// API v1
	r.Route("/api/v1", func(r chi.Router) {
		// Tasks endpoints
		taskHandler := handlers.NewTaskHandler(taskService, log)

		r.Route("/tasks", func(r chi.Router) {
			r.Post("/", taskHandler.CreateTask)
			r.Get("/", taskHandler.ListTasks)
			r.Get("/{id}", taskHandler.GetTask)
			r.Patch("/{id}/done", taskHandler.MarkTaskDone)
			r.Delete("/{id}", taskHandler.DeleteTask)
		})
	})

	return r
}

// LoggerMiddleware - кастомный middleware для логирования запросов
func LoggerMiddleware(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// оборачивание ResponseWriter для захвата статус кода
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			next.ServeHTTP(ww, r)

			duration := time.Since(start)

			log.Info("request completed",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status", ww.Status()),
				slog.Int("bytes", ww.BytesWritten()),
				slog.String("duration", duration.String()),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("request_id", middleware.GetReqID(r.Context())),
			)
		})
	}
}
