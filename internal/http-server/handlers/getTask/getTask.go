package getTask

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	resp "github.com/zeze322/todoapp/internal/lib/api/response"
	"github.com/zeze322/todoapp/internal/logger/sl"
	"github.com/zeze322/todoapp/internal/storage"
	"github.com/zeze322/todoapp/internal/storage/postgres"
	"log/slog"
	"net/http"
	"strconv"
)

type TaskGetter interface {
	GetTask(id int) (postgres.Task, error)
}

func New(log *slog.Logger, getTask TaskGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.getTask.New"

		log = slog.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		id := chi.URLParam(r, "id")
		if id == "" {
			log.Info("id not provided")

			render.JSON(w, r, resp.Error("invalid request"))

			return
		}

		idValue, err := strconv.Atoi(id)
		if err != nil {
			log.Info("bad id", id)

			render.JSON(w, r, resp.Error("bad id"))

			return
		}

		task, err := getTask.GetTask(idValue)
		if errors.Is(err, storage.ErrIDNotFound) {
			log.Info("task not found")

			render.JSON(w, r, resp.Error("not found"))

			return
		}

		if err != nil {
			log.Error("failed to get task", sl.Err(err))

			render.JSON(w, r, resp.Error("internal server error"))

			return
		}

		log.Info("got task")

		render.JSON(w, r, task)
	}
}
