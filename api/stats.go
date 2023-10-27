package api

import (
    "net/http"
    "encoding/json"
)

func (h *Handler) GetStats(w http.ResponseWriter, r *http.Request) {
    stats := h.s.Test()
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(stats)
}
