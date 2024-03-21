package update

import (
	h "awesomeProject/internal/http-server/handlers"
	resp "awesomeProject/internal/lib/api/response"
	"awesomeProject/internal/lib/logger/sl"
	"awesomeProject/internal/storage"
	"errors"
	"log/slog"
	"net/http"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias" validate:"required"`
}

type URLPUpdater interface {
	UpdateURL(alias string, newURL string) error
}

func New(log *slog.Logger, urlUpdater URLPUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.update.New"

		var req Request

		if err := h.HandleRequestBody(log, r, &req, w, op); err != nil {
			return
		}
		if err := h.ValidateRequestBody(log, r, &req, w, op); err != nil {
			return
		}

		err := urlUpdater.UpdateURL(req.Alias, req.URL)

		if err != nil {
			if errors.Is(err, storage.ErrAliasNotFound) {
				log.Info("alias not found", slog.String("alias", req.Alias))
				h.ResponseOK(w, r, resp.Error("alias not found"))
				return
			}

			log.Error("failed to update url", sl.Err(err))
			h.ResponseServerError(w, r, "failed to update url")

			return
		}

		log.Info("url updated", slog.String("url", req.URL))

		h.ResponseOK(w, r, resp.OK())
	}
}
