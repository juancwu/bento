package oauth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/juancwu/bento/env"
)

func generateRandomString(n int) (string, error) {
	data := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, data); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(data), nil
}

// generates a random state to use to identify the oauth redirect uri
func createOAuthState(secret string) (string, error) {
	randString, err := generateRandomString(32)
	if err != nil {
		return "", err
	}

	signature, err := hash(randString + secret)
	if err != nil {
		return "", err
	}

	stateString := randString + "." + signature

	return stateString, nil
}

func verifyState(state string, secret string) error {
	parts := strings.Split(state, ".")

	if len(parts) < 2 {
		return fmt.Errorf("OAUTH state string is invalid.")
	}

	temptableSignature := parts[1]
	randomString := parts[0]
	trueSignature, err := hash(randomString + secret)
	if err != nil {
		return err
	}

	if temptableSignature != trueSignature {
		return fmt.Errorf("OAUTH State signature does not match.")
	}

	return nil
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

func hash(s string) (string, error) {
	hasher := sha256.New()
	if _, err := io.WriteString(hasher, s); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(hasher.Sum(nil)), nil
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

func getUserInfoFromGitHub(token string) (*User, error) {
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

	var user User
	if err := json.NewDecoder(res.Body).Decode(&user); err != nil {
		return nil, errors.New("Error decoding response json from user info request")
	}

	return &user, nil
}

func createJWT(id int) (string, error) {
	now := time.Now()
	claims := OAuthToken{
		id,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(OAUTH_TOKEN_EXP)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv(env.SECRET_KEY)))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func verifyJWT(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &OAuthToken{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv(env.SECRET_KEY)), nil
	})

	if err != nil {
		return nil, err
	}

	if token.Valid {
		return token, nil
	}

	return nil, errors.New("Invalid token")
}
