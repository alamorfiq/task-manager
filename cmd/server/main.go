package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"taskmgr/internal/config"
	"taskmgr/internal/httpserver"
	"taskmgr/internal/service"
	"taskmgr/internal/storage/postgres"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg)

	log := setupLogger(cfg.Env)
	log.Info("starting task manager", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	// Инициализация БД
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := postgres.New(ctx, cfg.DB)
	if err != nil {
		log.Error("failed to init postgres", slog.String("err", err.Error()))
		os.Exit(1)
	}
	defer db.Close()

	log.Info("postgres connected successfully")

	// Инициализация сервиса
	taskService := service.NewTaskService(db)

	// Инициализация роутера
	router := httpserver.NewRouter(taskService, log)

	// Настройка HTTP сервера
	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	// Запуск сервера в горутине
	go func() {
		log.Info("starting http server", slog.String("address", cfg.HTTPServer.Address))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("failed to start server", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	log.Info("server started successfully")

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Info("shutting down server...")

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("server forced to shutdown", slog.String("error", err.Error()))
		os.Exit(1)
	}

	log.Info("server stopped gracefully")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
