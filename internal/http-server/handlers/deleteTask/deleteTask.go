package deleteTask

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	resp "github.com/zeze322/todoapp/internal/lib/api/response"
	"github.com/zeze322/todoapp/internal/logger/sl"
	"github.com/zeze322/todoapp/internal/storage"
)

type TaskDeleter interface {
	Delete(id int) error
}

func New(log *slog.Logger, deleteTask TaskDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.deleteTask.New"

		log = slog.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		id := chi.URLParam(r, "id")
		if id == "" {
			log.Info("id not provided")

			render.JSON(w, r, resp.Error("id not provided"))

			return
		}

		idValue, err := strconv.Atoi(id)
		if err != nil {
			log.Info("bad id", id)

			render.JSON(w, r, resp.Error("bad id"))

			return
		}

		err = deleteTask.Delete(idValue)
		if errors.Is(err, storage.ErrIDNotFound) {
			log.Info("task not found")

			render.JSON(w, r, resp.Error("task not found"))

			return
		}

		if err != nil {
			log.Error("failed to delete task", sl.Err(err))

			render.JSON(w, r, resp.Error("internal server error"))

			return
		}

		log.Info("task deleted")
	}
}
