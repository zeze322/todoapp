package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/zeze322/todoapp/internal/config"
	"github.com/zeze322/todoapp/internal/http-server/handlers/createTask"
	"github.com/zeze322/todoapp/internal/http-server/handlers/deleteTask"
	"github.com/zeze322/todoapp/internal/http-server/handlers/getTask"
	"github.com/zeze322/todoapp/internal/http-server/handlers/getTasks"
	"github.com/zeze322/todoapp/internal/http-server/handlers/updateTask"
	"github.com/zeze322/todoapp/internal/logger/sl"
	"github.com/zeze322/todoapp/internal/storage/postgres"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.Load()

	log := setupLogger(cfg.Env)

	storage, err := postgres.New()
	if err != nil {
		log.Error("could not connect to database", sl.Err(err))
		os.Exit(1)
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)

	router.Get("/task/{id}", getTask.New(log, storage))
	router.Get("/tasks", getTasks.New(log, storage))
	router.Post("/tasks", createTask.New(log, storage))
	router.Put("/task/{id}", updateTask.New(log, storage))
	router.Delete("/task/{id}", deleteTask.New(log, storage))

	srv := &http.Server{
		Addr:              cfg.Address,
		Handler:           router,
		ReadHeaderTimeout: cfg.Timeout,
		WriteTimeout:      cfg.Timeout,
		IdleTimeout:       cfg.IdleTimeout,
	}

	log.Info("starting server", slog.String("addr", srv.Addr))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err = srv.ListenAndServe(); err != nil {
			log.Error("failed starting server")
		}
	}()

	log.Info("server started")

	<-done

	log.Info("shutting down server")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(
			os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(slog.NewTextHandler(
			os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(slog.NewTextHandler(
			os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
