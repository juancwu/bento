package web

import (
	"github.com/go-chi/chi/v5"

	"github.com/juancwu/bento/store"
)

type Handler struct {
    router chi.Router
    store *store.Store
}
