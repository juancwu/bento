package oauth

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"

	"github.com/juancwu/bento/store"
)

type Handler struct {
    router chi.Router
    store *store.Store
}

func New(s *store.Store) *Handler {
    h := &Handler{}

    h.router = chi.NewRouter()

    h.router.Get("/state", h.GetRandomState)
    h.router.Get("/validate-state", h.VerifyOAuthState)

    h.store = s

    return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    h.router.ServeHTTP(w, r)
}

func (h *Handler) urlParam(r *http.Request, key string) string {
    return chi.URLParam(r, key)
}

func (h *Handler) GetRandomState(w http.ResponseWriter, r *http.Request) {
    state, err := CreateOAuthState(os.Getenv("SECRET_KEY"))
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("Error generating random state"))
    }
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(state))
}

func (h *Handler) VerifyOAuthState(w http.ResponseWriter, r *http.Request) {
    state := r.URL.Query().Get("state");

    if len(state) == 0 {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("No state provided to verify"))
    }

    err := VerifyState(state, os.Getenv("SECRET_KEY"))
    w.WriteHeader(http.StatusOK)
    if err != nil {
        w.Write([]byte("Invalid state"))
    } else {
        w.Write([]byte("Valid state"))
    }
}
