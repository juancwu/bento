package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/juancwu/bento/store"
)

type Handler struct {
    router chi.Router
    store *store.Store
}

func New(store *store.Store) *Handler {
    h := &Handler{}

    h.router = chi.NewRouter()

    h.router.Get("/bentos", h.GetBentos)
    h.router.Get("/stats", h.GetStats)

    h.store = store

    return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    h.router.ServeHTTP(w, r)
}

func (h *Handler) urlParam(r *http.Request, key string) string {
    return chi.URLParam(r, key)
}
