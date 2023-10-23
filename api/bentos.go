package api

import (
    "net/http"
)

func (h *Handler) GetBentos(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("from bentos api"))
}
