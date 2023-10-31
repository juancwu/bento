package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"

    "github.com/juancwu/bento/web"
	"github.com/juancwu/bento/store"
)

func New(store *store.Store) *web.Handler {
    h := &web.Handler{}

    h.router = chi.NewRouter()

    h.router.Get("/bentos", h.GetBentos)
    h.router.Get("/stats", h.GetStats)

    h.store = store

    return h
}

func (h *web.Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    h.router.ServeHTTP(w, r)
}

func (h *web.Handler) urlParam(r *http.Request, key string) string {
    return chi.URLParam(r, key)
}
