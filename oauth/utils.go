package oauth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
    "encoding/json"
	"fmt"
	"io"
    "os"
	"strings"
    "net/url"
    "net/http"
    "errors"

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
