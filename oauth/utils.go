package oauth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"strings"
)

func GenerateRandomString(n int) (string, error) {
    data := make([]byte, n)
    if _, err := io.ReadFull(rand.Reader, data); err != nil {
        return "", err
    }

    return base64.StdEncoding.EncodeToString(data), nil
}

// generates a random state to use to identify the oauth redirect uri
func CreateOAuthState(secret string) (string, error) {
    randString, err := GenerateRandomString(32)
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

func VerifyState(state string, secret string) error {
    parts := strings.Split(state, ".")

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

func hash(s string) (string, error) {
    hasher := sha256.New()
    if _, err := io.WriteString(hasher, s); err != nil {
        return "", err
    }

    return base64.StdEncoding.EncodeToString(hasher.Sum(nil)), nil
}
