package oauth

import (
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"

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
    Id int `json:"id"`
    Email *string `json:"email"`
    Login *string `json:"login"`
}

type OAuthSuccessResponse struct {
    Token string `json:"token"`
}

type OAuthTokenJWT struct {
    Id int `json:"id"`
    jwt.RegisteredClaims
}

type OAuthStateJWT struct {
    State string `json:"state"`
    Port string `json:"port"`
    Cli bool `json:"cli"`
    jwt.RegisteredClaims
}
