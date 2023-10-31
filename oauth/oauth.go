package oauth

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"net/http"

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

    h.store = s

    return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    h.router.ServeHTTP(w, r)
}

func (h *Handler) urlParam(r *http.Request, key string) string {
    return chi.URLParam(r, key)
}

// Generates a random state to use to identify the oauth redirect uri
func State(n int) (string, error) {
    data := make([]byte, n)
    if _, err := io.ReadFull(rand.Reader, data); err != nil {
        return "", err
    }

    return base64.StdEncoding.EncodeToString(data), nil
}

