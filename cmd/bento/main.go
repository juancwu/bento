package main

import (
	"fmt"
    "net/http"
    "time"
    "encoding/json"
    "log"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
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

    r.Get("/", handler)

    r.Route("/grouped", func(r chi.Router) {
        r.Get("/", func(w http.ResponseWriter, r *http.Request) {
            w.Write([]byte("base route for grouped routes"))
        })

        r.Post("/", func(w http.ResponseWriter, r *http.Request) {
            var body PostBody

            err := json.NewDecoder(r.Body).Decode(&body)
            if err != nil {
                http.Error(w, "Failed to decode JSON body", http.StatusBadRequest)
                return
            }

            log.Printf("Received: %+v", body)

            w.WriteHeader(http.StatusOK)
            w.Write([]byte("Data received!"))
        })
    })

    addr := ":3000"
    fmt.Printf("Serving on port %s\n", addr)
    http.ListenAndServe(addr, r)
}

func handler(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Welcome"))
}
