package getTasks

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	resp "github.com/zeze322/todoapp/internal/lib/api/response"
	"github.com/zeze322/todoapp/internal/logger/sl"
	"github.com/zeze322/todoapp/internal/storage"
	"github.com/zeze322/todoapp/internal/storage/postgres"
	"log/slog"
	"net/http"
)

type TasksGetter interface {
	Tasks() ([]postgres.Task, error)
}

func New(log *slog.Logger, getTasks TasksGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.getTasks.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		tasks, err := getTasks.Tasks()
		if errors.Is(err, storage.ErrTasksNotFound) {
			log.Info("tasks not found")

			render.JSON(w, r, resp.Error("tasks not found"))

			return
		}

		if err != nil {
			log.Error("failed to get tasks", sl.Err(err))

			render.JSON(w, r, "internal server error")

			return
		}

		log.Info("got tasks")

		render.JSON(w, r, tasks)
	}
}
