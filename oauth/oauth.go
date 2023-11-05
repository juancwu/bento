package oauth

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/go-chi/chi/v5"

	"github.com/juancwu/bento/env"
	"github.com/juancwu/bento/store"
	"github.com/juancwu/bento/web"
)

const (
	GITHUB_OAUTH_URL       = "https://github.com/login/oauth/authorize"
	GITHUB_TOKEN_URL       = "https://github.com/login/oauth/access_token"
	GITHUB_USER_EMAILS_URL = "https://api.github.com/user/emails"
	GITHUB_USER_INFO_URL   = "https://api.github.com/user"
	NANOID_LEN             = 12
)

func New(s *store.Store) *OAuthHandler {
	h := &OAuthHandler{}

	h.router = chi.NewRouter()

	h.router.Get("/github/callback", h.HandleCallback)
	h.router.Get("/login", h.GetLoginPage)

	h.store = s

	return h
}

func (h *OAuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

func (h *OAuthHandler) urlParam(r *http.Request, key string) string {
	return chi.URLParam(r, key)
}

func (h *OAuthHandler) GetLoginPage(w http.ResponseWriter, r *http.Request) {
	query := url.Values{}
	query.Set("client_id", os.Getenv(env.GITHUB_OAUTH_CLIENT_ID))
	query.Set("scope", "user:email")
	state, err := createOAuthState(os.Getenv(env.SECRET_KEY))
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Problem occurred getting url"))
		return
	}
	query.Set("state", state)
	authURL := fmt.Sprintf("%s?%s", GITHUB_OAUTH_URL, query.Encode())
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

func (h *OAuthHandler) HandleCallback(w http.ResponseWriter, r *http.Request) {
	// verify state
	state := r.URL.Query().Get("state")
	if state == "" {
		fmt.Println("OAuth callback attempt without state")
		http.Error(w, "No state", http.StatusBadRequest)
		return
	}

	err := verifyState(state, os.Getenv(env.SECRET_KEY))
	if err != nil {
		fmt.Println("OAuth callback attempt with invalid state")
		http.Error(w, "Invalid state", http.StatusBadRequest)
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		fmt.Println("OAuth callback attempt without code")
		http.Error(w, "No code", http.StatusBadRequest)
		return
	}

	token, err := exchangeCodeForToken(code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userInfo, err := getUserInfoFromGitHub(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if userInfo.Email == nil {
		// create new user if not exists
		// get user email, and it has to be verified
		email, err := getUserPrimaryEmail(token)
		if err != nil {
			msg := err.Error()
			fmt.Println(msg)
			http.Error(w, msg, http.StatusBadRequest)
		}

		fmt.Printf("User Email: %s\n", email)

		userInfo.Email = &email
	}

	user, err := h.store.CreateNewUser(*userInfo.Email, userInfo.ID)
	if err != nil {
		fmt.Printf("ERROR: could not create new user: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := OAuthSuccessResponse{}
	response.Email = user.Email
	response.ObjectId = user.ObjectId
	response.GhId = userInfo.ID
	response.Token = token

	// TODO: create signed jwt to send back to client

	web.Json(w, response, http.StatusOK)
}
