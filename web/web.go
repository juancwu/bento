package web

import (
	"encoding/json"
	"net/http"
)

func Reply(w http.ResponseWriter, message string, code int) {
    w.WriteHeader(code)
    w.Write([]byte(message))
}

func Json(w http.ResponseWriter, data interface{}, code int) {
    w.Header().Set("Content-Type", "application/json")

    jsonData, err := json.Marshal(data)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte(`{"error": "Something went wrong marshalling the response object"}`))
        return
    }

    w.WriteHeader(code)
    w.Write(jsonData)
}
