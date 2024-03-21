package delete

import (
	h "awesomeProject/internal/http-server/handlers"
	resp "awesomeProject/internal/lib/api/response"
	"awesomeProject/internal/lib/logger/sl"
	"awesomeProject/internal/storage"
	"errors"
	"log/slog"
	"net/http"
)

type URLDelete interface {
	DeleteURL(alias string) error
}

func New(log *slog.Logger, urlDelete URLDelete) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.delete.New"

		params, err := h.HandleRequestParams(log, r, w, op, "alias")
		if err != nil {
			return
		}

		err = urlDelete.DeleteURL(params["alias"])
		if err != nil {
			if errors.Is(err, storage.ErrAliasNotFound) {
				log.Info("alias not found", slog.String("alias", params["alias"]))
				h.ResponseOK(w, r, resp.Error("alias not found"))
				return
			}

			log.Error("failed to delete alias", sl.Err(err))
			h.ResponseServerError(w, r, "failed to delete alias")

			return
		}

		log.Info("alias deleted", slog.String("alias", params["alias"]))

		h.ResponseOK(w, r, resp.OK())
	}
}
