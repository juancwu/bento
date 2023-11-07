package oauth

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

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
	OAUTH_TOKEN_EXP        = 168 * time.Hour // one week 
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
	// check for bearer token
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		parts := strings.Fields(authHeader)
		if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
			tokenString := parts[1]
			fmt.Printf("Token: %s\n", tokenString)
			// verify the token
			token, err := verifyJWT(tokenString)
			if err != nil {
				fmt.Printf("ERROR: %s\n", err.Error())
			} else {
				claims, ok := token.Claims.(*OAuthTokenJWT)
				if ok {
					fmt.Printf("Claims: %v\n", claims)
					// get user info
					user, err := h.store.GetUserById(claims.Id)
					if err != nil {
						fmt.Printf("Could not get user with id: %d\n", claims.Id)
						fmt.Printf("ERROR: %s\n", err.Error())
					} else {
						fmt.Printf("User: %v\n", user)
					}
				}

				fmt.Fprint(w, "Token is valid")
				return
			}
		}
	}

    flow := strings.ToLower(r.URL.Query().Get("flow"))
    redirect := strings.ToLower(r.URL.Query().Get("redirect"))
    portString := r.URL.Query().Get("port")

    if redirect == "" || strings.Contains(redirect, "localhost") {
        http.Error(w, "Invalid redirect query parameter", http.StatusBadRequest)
        return
    }

    if flow != "" && flow != "web" && flow != "cli" {
        http.Error(w, "Invalid flow query parameter", http.StatusBadRequest)
        return
    }

    validPort, port := isValidPort(portString)

    if !validPort {
        http.Error(w, "Invalid port query parameter", http.StatusBadRequest)
        return
    }

	query := url.Values{}
	query.Set("client_id", os.Getenv(env.GITHUB_OAUTH_CLIENT_ID))
	query.Set("scope", "user:email")

	state, stateId, err := createOAuthState()
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Problem occurred getting url"))
		return
	}

    // save state in database for further query on callback
    err = h.store.SaveState(stateId, flow, redirect, port)
    if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Problem saving state"))
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

	ghUser, err := getUserInfoFromGitHub(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if ghUser.Email == nil {
		// create new user if not exists
		// get user email, and it has to be verified
		email, err := getUserPrimaryEmail(token)
		if err != nil {
			msg := err.Error()
			fmt.Println(msg)
			http.Error(w, msg, http.StatusBadRequest)
		}

		fmt.Printf("User Email: %s\n", email)

		ghUser.Email = &email
	}

	// check if user exists
	user, err := h.store.GetUserByGhId(ghUser.Id)
	if err != nil && err == sql.ErrNoRows {
		user, err = h.store.CreateNewUser(*ghUser.Email, ghUser.Id)
		if err != nil {
			fmt.Printf("ERROR: could not create new user: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

    ghToken, err := h.store.GetAccessToken(user.Id)
    if err == sql.ErrNoRows {
        // create new access token entry
        h.store.SaveAccessToken(user.Id, token)
    } else if err != nil {
        fmt.Printf("ERROR: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    } else if  ghToken != token {
        // update acess token in db
        err = h.store.UpdateAccessToken(user.Id, token)
        if err != nil {
            fmt.Printf("ERROR: %v", err)
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
    }
    claims := OAuthTokenJWT{
        Id: user.Id,
        RegisteredClaims: getStdJWTClaims(OAUTH_TOKEN_EXP),
    }

	tokenString, err := createJWT(claims)
	if err != nil {
		fmt.Printf("ERROR: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := OAuthSuccessResponse{
		Token: tokenString,
	}

	web.Json(w, response, http.StatusOK)
}
