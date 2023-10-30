package api

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"

	"github.com/juancwu/bento/store"
)

type Handler struct {
    router chi.Router
    store *store.Store
}

func New() *Handler {
    h := &Handler{}

    h.router = chi.NewRouter()

    h.router.Get("/bentos", h.GetBentos)
    h.router.Get("/stats", h.GetStats)

    s, err := store.New()
    if err != nil {
        fmt.Fprintf(os.Stderr, "failed to open db %s: %s", os.Getenv("BENTO_DB_URL"), err)
        os.Exit(1)
    }

    h.store = s

    return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    h.router.ServeHTTP(w, r)
}

func (h *Handler) urlParam(r *http.Request, key string) string {
    return chi.URLParam(r, key)
}
