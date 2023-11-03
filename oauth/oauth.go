package oauth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"

	"github.com/juancwu/bento/store"
)

const GITHUB_OAUTH_URL = "https://github.com/login/oauth/authorize?scope=user:email&state=%s&client_id=%s"
const GITHUB_OAUTH_ACCESS_URL = "https://github.com/login/oauth/access_token?clieant_id=%s&client_secret=%s&code=%s"

type Handler struct {
    router chi.Router
    store *store.Store
}

type OAuthAcessRes struct {
    accessToken string `json:"access_token"`
    scope string `json:"scope"`
    tokenType string `json:"token_type"`
}

func New(s *store.Store) *Handler {
    h := &Handler{}

    h.router = chi.NewRouter()

    // TODO: remove this routes once done with oauth implementation
    h.router.Get("/state", h.GetRandomState)
    h.router.Get("/validate-state", h.VerifyOAuthState)
    h.router.Get("/url", h.GetOAuthURL)
    h.router.Get("/github/callback", h.HandleOAuthGrant)

    h.store = s

    return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    h.router.ServeHTTP(w, r)
}

func (h *Handler) urlParam(r *http.Request, key string) string {
    return chi.URLParam(r, key)
}

func (h *Handler) GetRandomState(w http.ResponseWriter, r *http.Request) {
    state, err := CreateOAuthState(os.Getenv("SECRET_KEY"))
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("Error generating random state"))
    }
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(state))
}

func (h *Handler) VerifyOAuthState(w http.ResponseWriter, r *http.Request) {
    state := r.URL.Query().Get("state");

    if len(state) == 0 {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("No state provided to verify"))
    }

    err := VerifyState(state, os.Getenv("SECRET_KEY"))
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("Invalid state"))
    } else {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("Valid state"))
    }
}

func (h *Handler) GetOAuthURL(w http.ResponseWriter, r *http.Request) {
    state, err := CreateOAuthState(os.Getenv("SECRET_KEY"))
    if err != nil {
        fmt.Printf("ERROR: %v\n", err)
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("Problem occurred getting url"))
    }

    url := GetGitHubOAuthURL(state)

    w.WriteHeader(http.StatusOK)
    w.Write([]byte(url))
}

func (h *Handler) HandleOAuthGrant(w http.ResponseWriter, r *http.Request) {
    code := r.URL.Query().Get("code")
    state := r.URL.Query().Get("state")
    if len(code) == 0 || len(state) == 0 {
        // not granted
        w.WriteHeader(http.StatusUnauthorized)
        w.Write([]byte("OAuth not granted"))
    }

    // verify state
    err := VerifyState(state, os.Getenv("SECRET_ENV"))
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("Invalid state"))
        return
    }

    // get access token
    accessURL := fmt.Sprintf(GITHUB_OAUTH_ACCESS_URL, os.Getenv("OAUTH_GITHUB_CLIENT_ID"), os.Getenv("OAUTH_GITHUB_CLIENT_SECRET"), code)
    req, err := http.NewRequest(http.MethodGet, accessURL, nil)
    if err != nil {
        fmt.Printf("ERROR: could not initialize request: %s", err)
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte("Could not initialize request to get access token"))
        return
    }

    req.Header.Set("Accept", "application/json")
    res, err := http.Get(accessURL)
    if err != nil {
        fmt.Printf("ERROR: could not complete http request to get access token: %s", err)
        w.WriteHeader(res.StatusCode)
        w.Write([]byte("Unable to complete OAuth flow"))
        return
    }
    defer res.Body.Close()

    if res.StatusCode != http.StatusOK {
        fmt.Printf("ERROR: Bad OAuth http request: %s", err)
        w.WriteHeader(res.StatusCode)
        w.Write([]byte("Bad OAuth http request"))
    }

    var oauthAccessRes OAuthAcessRes
    err = json.NewDecoder(res.Body).Decode(&oauthAccessRes)
    if err != nil {
        fmt.Printf("ERROR: could not parse response body: %s", err)
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte("Could not parse response body"))
    }

    fmt.Printf("Access Token: %s\nScope: %s\nToken Type: %s\n", oauthAccessRes.accessToken, oauthAccessRes.scope, oauthAccessRes.tokenType)

    w.WriteHeader(http.StatusOK)
    w.Header().Set("Content-Type", "application/json")

    _, err = io.Copy(w, res.Body)
    if err != nil {
        fmt.Printf("ERROR: could not send json response to client: %v", err)
    }
}
