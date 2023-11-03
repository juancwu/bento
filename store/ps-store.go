package store

import (
	"database/sql"
	"os"

	"github.com/juancwu/bento/env"
)

func NewPsStore() (*Store, error) {
    db, err := sql.Open("mysql", os.Getenv(env.DSN))
    if err != nil {
        return nil, err
    }

    s := &Store{}
    s.db = db

    return s, nil
}
