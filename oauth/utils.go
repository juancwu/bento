package oauth

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	gonanoid "github.com/matoous/go-nanoid/v2"

	"github.com/juancwu/bento/env"
)

// generates a random state to use to identify the oauth redirect uri
func createOAuthState(cli bool, port string) (string, error) {
	randString, err := gonanoid.New(32)
	if err != nil {
		return "", err
	}

	signature, err := hash(randString + os.Getenv(env.SECRET_KEY))
	if err != nil {
		return "", err
	}

	// create jwt to send as state
	jwt := OAuthStateJWT{
		Signature:        signature,
		State:            randString,
		Port:             port,
		Cli:              cli,
		RegisteredClaims: getStdJWTClaims(10 * time.Minute),
	}
	jwtString, err := createJWT(jwt)
	if err != nil {
		return "", err
	}

	return jwtString, nil
}

func verifyOAuthState(stateJWT string) (*OAuthStateJWT, error) {
	token, err := jwt.ParseWithClaims(stateJWT, &OAuthStateJWT{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv(env.SECRET_KEY)), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("Invalid token")
	}

	jwt, ok := token.Claims.(*OAuthStateJWT)
	if !ok {
		return nil, errors.New("Invalid JWT format")
	}

	trueSignature, err := hash(jwt.State + os.Getenv(env.SECRET_KEY))
	if err != nil {
		return nil, err
	}

	if jwt.Signature != trueSignature {
		return nil, errors.New("OAUTH State signature does not match.")
	}

	return jwt, nil
}

func hash(s string) (string, error) {
	hasher := sha256.New()
	if _, err := io.WriteString(hasher, s); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(hasher.Sum(nil)), nil
}

func exchangeCodeForToken(code string) (string, error) {
	query := url.Values{}
	query.Set("client_id", os.Getenv(env.GITHUB_OAUTH_CLIENT_ID))
	query.Set("client_secret", os.Getenv(env.GITHUB_OAUTH_CLIENT_SECRET))
	query.Set("code", code)

	req, err := http.NewRequest(http.MethodPost, GITHUB_TOKEN_URL, strings.NewReader(query.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var data map[string]interface{}
	if err = json.NewDecoder(res.Body).Decode(&data); err != nil {
		return "", err
	}

	token, ok := data["access_token"].(string)
	if !ok {
		return "", errors.New("Invalid token response")
	}

	return token, nil
}

func getUserPrimaryEmail(token string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, GITHUB_USER_EMAILS_URL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", errors.New(fmt.Sprintf("Server returned non-200 status: %d %s", res.StatusCode, res.Status))
	}

	var emails []Email
	if err := json.NewDecoder(res.Body).Decode(&emails); err != nil {
		return "", errors.New("Error decoding response json from user emails request")
	}

	for _, email := range emails {
		if email.Primary {
			if !email.Verified {
				return "", errors.New("OAuth callback attempt, user email not verified")
			}

			return email.Email, nil
		}
	}

	return "", errors.New("No email available")
}

func getUserInfoFromGitHub(token string) (*GitHubUser, error) {
	req, err := http.NewRequest(http.MethodGet, GITHUB_USER_INFO_URL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Server returned non-200 status: %s", res.Status))
	}

	var user GitHubUser
	if err := json.NewDecoder(res.Body).Decode(&user); err != nil {
		return nil, errors.New("Error decoding response json from user info request")
	}

	return &user, nil
}

func createJWT(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv(env.SECRET_KEY)))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func keyFuncJWT(token *jwt.Token) (interface{}, error) {
	_, ok := token.Method.(*jwt.SigningMethodHMAC)
	if !ok {
		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	}

	return []byte(os.Getenv(env.SECRET_KEY)), nil
}

func getStdJWTClaims(exp time.Duration) jwt.RegisteredClaims {
	now := time.Now()
	stdClaims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(now.Add(exp)),
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
	}

	return stdClaims
}

func isValidPort(portStr string) bool {
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return false
	}

	valid := port > 0 && port <= 65535
	if !valid {
		return false
	}

	return true
}
