package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/juancwu/bento/api"
	"github.com/juancwu/bento/env"
	"github.com/juancwu/bento/oauth"
	"github.com/juancwu/bento/store"
)

type PostBody struct {
    Value string `json:"value"`
}

func main() {
    fmt.Println("Load env...")
    err := env.Load()
    if err != nil {
        panic(err)
    }

    r := chi.NewRouter()
    r.Use(middleware.RequestID)
    r.Use(middleware.RealIP)
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)

    r.Use(middleware.Timeout(60 * time.Second))

    r.Get("/", func (w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("It works!"))
    })

    s, err := store.New()
    if err != nil {
        fmt.Fprintf(os.Stderr, "failed to open db %s: %s", os.Getenv("BENTO_DB_URL"), err)
        os.Exit(1)
    }

    apiHandler := api.New(s)
    oauthHandler := oauth.New(s)
    r.Mount("/api/v1", apiHandler)
    r.Mount("/oauth", oauthHandler)

    addr := ":3000"
    fmt.Printf("Serving on port %s\n", addr)
    http.ListenAndServe(addr, r)
}
