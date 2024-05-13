package updateTask

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	resp "github.com/zeze322/todoapp/internal/lib/api/response"
	"github.com/zeze322/todoapp/internal/logger/sl"
)

type Request struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type TaskUpdater interface {
	Update(id int, title, description string) error
}

func New(log *slog.Logger, updateTask TaskUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handle.updateTask.New"

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

		var req Request

		err = render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("empty request body")

			render.JSON(w, r, resp.Error("empty request body"))

			return
		}

		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		err = updateTask.Update(idValue, req.Title, req.Description)
		if err != nil {
			log.Error("failed to update task", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to update task"))

			return
		}

		log.Info("task updated")
	}
}
