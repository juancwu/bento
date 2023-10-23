package main

import (
	"fmt"
    "net/http"
    "time"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"

    "github.com/juancwu/bento/api"
)

type PostBody struct {
    Value string `json:"value"`
}

func main() {
    r := chi.NewRouter()
    r.Use(middleware.RequestID)
    r.Use(middleware.RealIP)
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)

    r.Use(middleware.Timeout(60 * time.Second))

    r.Get("/", func (w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("It works!"))
    })

    apiHandler := api.New()
    r.Mount("/api", apiHandler)

    addr := ":3000"
    fmt.Printf("Serving on port %s\n", addr)
    http.ListenAndServe(addr, r)
}
