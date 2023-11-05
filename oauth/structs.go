package oauth

import (
	"github.com/go-chi/chi/v5"

	"github.com/juancwu/bento/store"
)

type OAuthHandler struct {
	router chi.Router
	store  *store.Store
}

type Email struct {
    Email string `json:"email"`
    Primary bool `json:"primary"`
    Verified bool `json:"verified"`
    Visibility *string `json:"visibility"`
}

type User struct {
    ID int `json:"id"`
    Email *string `json:"email"`
    Login *string `json:"login"`
}

type OAuthSuccessResponse struct {
    Token string `json:"token"`
    Email string `json:"email"`
    ObjectId string `json:"object_id"`
    GhId int `json:"gh_id"`
}
