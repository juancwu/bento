package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
    router chi.Router
}

func New() *Handler {
    h := &Handler{}

    h.router = chi.NewRouter()

    h.router.Get("/bentos", h.GetBentos)

    return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    h.router.ServeHTTP(w, r)
}

func (h *Handler) urlParam(r *http.Request, key string) string {
    return chi.URLParam(r, key)
}
