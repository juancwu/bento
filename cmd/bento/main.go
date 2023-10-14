package main

import (
	"fmt"
    "net/http"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
)

func main() {
    r := chi.NewRouter()
    r.Use(middleware.Logger)

    r.Get("/", handler)

    addr := ":3000"
    fmt.Printf("Serving on port %s\n", addr)
    http.ListenAndServe(addr, r)
}

func handler(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Welcome"))
}
