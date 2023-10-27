package store

import (
    "database/sql"
    "os"

    _ "github.com/libsql/libsql-client-go/libsql"
)

var dbUrl = os.Getenv("BENTO_DB_URL") + "?authToken=" + os.Getenv("BENTO_DB_AUTH_TOKEN")

type Store struct {
    db *sql.DB
}

func New() (*Store, error) {
    db, err := sql.Open("libsql", dbUrl)
    if err != nil {
        return nil, err
    }

    s := &Store{}
    s.db = db

    return s, nil
}

func (s *Store) Test() sql.DBStats {
    return s.db.Stats()
}
