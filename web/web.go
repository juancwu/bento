package web

import "net/http"

func Reply(w http.ResponseWriter, message string, code int) {
    w.WriteHeader(code)
    w.Write([]byte(message))
}
