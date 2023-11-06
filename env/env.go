package env

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

const (
	SECRET_KEY                 = "SECRET_KEY"
	STATE_SECRET_KEY           = "STATE_SECRET_KEY"
	BENTO_DB_URL               = "BENTO_DB_URL"
	BENTO_DB_AUTH_TOKEN        = "BENTO_DB_AUTH_TOKEN"
	BENTO_DB_CONN              = "BENTO_DB_CONN"
	GITHUB_OAUTH_CLIENT_ID     = "GITHUB_OAUTH_CLIENT_ID"
	GITHUB_OAUTH_CLIENT_SECRET = "GITHUB_OAUTH_CLIENT_SECRET"
	DSN                        = "DSN"
)

func Load() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}

	err = checkEnv()
	if err != nil {
		return err
	}

	return nil
}

func checkEnv() error {
	envList := []string{
		SECRET_KEY,
        STATE_SECRET_KEY,
		BENTO_DB_URL,
		BENTO_DB_AUTH_TOKEN,
		BENTO_DB_CONN,
		GITHUB_OAUTH_CLIENT_ID,
		GITHUB_OAUTH_CLIENT_SECRET,
		DSN,
	}

	for _, name := range envList {
		value := os.Getenv(name)
		if len(value) == 0 {
			return fmt.Errorf("Missing mandatory env variable %s", name)
		}
	}

	return nil
}
