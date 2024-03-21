package redirect

import (
	h "awesomeProject/internal/http-server/handlers"
	"awesomeProject/internal/lib/logger/sl"
	"awesomeProject/internal/storage"
	"errors"
	"log/slog"
	"net/http"
)

type URLGetter interface {
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.redirect.New"

		params, err := h.HandleRequestParams(log, r, w, op, "alias")
		if err != nil {
			return
		}

		resURL, err := urlGetter.GetURL(params["alias"])
		if err != nil {
			if errors.Is(err, storage.ErrURLNotFound) {
				log.Info("url not found", "alias", params["alias"])
				h.ResponseClientError(w, r, "alias not found")
				return
			}

			log.Error("failed to get url", sl.Err(err))
			h.ResponseServerError(w, r, "failed to get url")
			return
		}

		log.Info("got url", slog.String("url", resURL))

		http.Redirect(w, r, resURL, http.StatusFound)
	}
}
