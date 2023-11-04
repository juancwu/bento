package store

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"github.com/juancwu/bento/env"
)

func NewPsStore() (*Store, error) {
    db, err := sql.Open("mysql", os.Getenv(env.DSN))
    if err != nil {
        return nil, err
    }

    err = db.Ping()
    if err != nil {
        return nil, err
    }

    fmt.Println("Sucessfully connected to PlanetScale DB")

    s := &Store{}
    s.db = db

    return s, nil
}
