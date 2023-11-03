package store

import (
	"database/sql"
	"fmt"
	"net/url"
	"os"

	"github.com/juancwu/bento/env"
	_ "github.com/libsql/libsql-client-go/libsql"
)

func NewTursoStore() (*Store, error) {
    query := url.Values{}
    query.Set("authToken", os.Getenv(env.BENTO_DB_AUTH_TOKEN))
    tursoStoreURL := fmt.Sprintf("%s?%s", os.Getenv(env.BENTO_DB_URL), query.Encode())
    db, err := sql.Open("libsql", tursoStoreURL)
    if err != nil {
        return nil, err
    }

    s := &Store{}
    s.db = db

    return s, nil
}
