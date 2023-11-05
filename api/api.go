package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/juancwu/bento/store"
	"github.com/juancwu/bento/web"
)

type Handler struct {
    router chi.Router
    store *store.Store
}

type Response struct {
    Message string `json:"message"`
}

func New(s *store.Store) *Handler {
    h := &Handler{}

    h.router = chi.NewRouter()

    h.router.Get("/bentos", h.GetBentos)
    h.router.Get("/stats", h.GetStats)
    h.router.Get("/json", func (w http.ResponseWriter, r *http.Request) {
        res := Response{
            Message: "hello, this is a json response",
        }
        web.Json(w, res, http.StatusOK)
    })

    h.store = s

    return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    h.router.ServeHTTP(w, r)
}

func (h *Handler) urlParam(r *http.Request, key string) string {
    return chi.URLParam(r, key)
}
