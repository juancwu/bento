package web

import (
    "fmt"
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

func Error(w http.ResponseWriter, err error, code int) {
    fmt.Printf("ERROR: %v", err)
    http.Error(w, err.Error(), code)
}
