package handlers

import (
	resp "awesomeProject/internal/lib/api/response"
	"awesomeProject/internal/lib/logger/sl"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
)

func HandleRequestBody(log *slog.Logger, r *http.Request, req interface{}, w http.ResponseWriter, op string) error {
	log = log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	err := render.DecodeJSON(r.Body, req)
	if err != nil {
		log.Error("failed to decode request body", sl.Err(err))
		ResponseServerError(w, r, "failed to decode request")
		return err
	}

	log.Info("request body decoded", slog.Any("request", req))
	return nil
}

func ValidateRequestBody(log *slog.Logger, r *http.Request, req interface{}, w http.ResponseWriter, op string) error {
	log = log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	if err := validator.New().Struct(req); err != nil {
		var validateErr validator.ValidationErrors
		errors.As(err, &validateErr)

		log.Error("invalid request", sl.Err(err))

		ResponseClientError(w, r, resp.ValidationError(validateErr))

		return err
	}

	return nil
}

func HandleRequestParams(log *slog.Logger, r *http.Request, w http.ResponseWriter, op string, keys ...string) (map[string]string, error) {
	log = log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)
	params := make(map[string]string)

	for _, key := range keys {
		value := chi.URLParam(r, key)
		if value == "" {
			err := fmt.Errorf("missing parameter: %s", key)
			log.Error("invalid request", err)
			ResponseClientError(w, r, fmt.Sprintf("missing parameter: %s", key))
			return nil, err
		}
		params[key] = value
	}

	return params, nil
}
