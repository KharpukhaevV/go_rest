package handlers

import (
	resp "awesomeProject/internal/lib/api/response"
	"github.com/go-chi/render"
	"net/http"
)

func ResponseOK(w http.ResponseWriter, r *http.Request, data interface{}) {
	render.JSON(w, r, data)
}

func ResponseServerError(w http.ResponseWriter, r *http.Request, errorMessage string) {
	w.WriteHeader(http.StatusInternalServerError)
	render.JSON(w, r, resp.Error(errorMessage))
}

func ResponseClientError(w http.ResponseWriter, r *http.Request, errorMessage string) {
	w.WriteHeader(http.StatusBadRequest)
	render.JSON(w, r, resp.Error(errorMessage))
}
