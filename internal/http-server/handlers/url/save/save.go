package save

import (
	h "awesomeProject/internal/http-server/handlers"
	resp "awesomeProject/internal/lib/api/response"
	"awesomeProject/internal/lib/logger/sl"
	"awesomeProject/internal/lib/random"
	"awesomeProject/internal/storage"
	"errors"
	"log/slog"
	"net/http"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias"`
}

type URLSaver interface {
	SaveUrl(urlToSave string, alias string) (int64, error)
}

const aliasLength = 6

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		var req Request

		if err := h.HandleRequestBody(log, r, &req, w, op); err != nil {
			return
		}
		if err := h.ValidateRequestBody(log, r, &req, w, op); err != nil {
			return
		}

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLength)
		}

		id, err := urlSaver.SaveUrl(req.URL, alias)

		if err != nil {
			if errors.Is(err, storage.ErrURLExists) {
				log.Info("url already exists", slog.String("url", req.URL))
				h.ResponseOK(w, r, resp.Error("url already exists"))
				return
			}

			log.Error("failed to add url", sl.Err(err))
			h.ResponseServerError(w, r, "failed to add url")

			return
		}

		log.Info("url added", slog.Int64("id", id))

		h.ResponseOK(w, r, Response{
			Response: resp.OK(),
			Alias:    alias,
		})
	}
}
