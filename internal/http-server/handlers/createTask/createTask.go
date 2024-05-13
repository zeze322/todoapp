package createTask

import (
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	resp "github.com/zeze322/todoapp/internal/lib/api/response"
	"github.com/zeze322/todoapp/internal/logger/sl"
	"github.com/zeze322/todoapp/internal/storage"
)

type Request struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type TaskCreator interface {
	Create(title, description string) error
}

func New(log *slog.Logger, createTask TaskCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.createTask.New"

		log = slog.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("empty request body")

			render.JSON(w, r, resp.Error("empty request"))

			return
		}

		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		err = createTask.Create(req.Title, req.Description)
		if errors.Is(err, storage.ErrTaskExists) {
			log.Info("task already exists", sl.Err(err))

			render.JSON(w, r, resp.Error("task already exists"))

			return
		}

		if err != nil {
			log.Error("failed to create task", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to create task"))

			return
		}

		log.Info("created task")
	}
}
