package env

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
    SECRET_KEY string
    BENTO_DB_URL string
    BENTO_DB_AUTH_TOKEN string
    BENTO_DB_CONN string
}

var EnvKey Env

func Load() error {
    err := godotenv.Load()
    if err != nil {
        return err
    }

    EnvKey = Env{
        SECRET_KEY: "SECRET_KEY",
        BENTO_DB_URL: "BENTO_DB_URL",
        BENTO_DB_AUTH_TOKEN: "BENTO_DB_AUTH_TOKEN",
        BENTO_DB_CONN: "BENTO_DB_CONN",
    }

    err = checkEnv()
    if err != nil {
        return err
    }

    return nil
}

func checkEnv() error {
    envList := []string{
        "SECRET_KEY",
        "BENTO_DB_URL",
        "BENTO_DB_AUTH_TOKEN",
        "BENTO_DB_CONN",
    }

    for _, name := range envList {
        value := os.Getenv(name)
        if len(value) == 0 {
            return fmt.Errorf("Missing mandatory env variable %s", name)
        }
    }

    return nil
}
