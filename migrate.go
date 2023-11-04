package main

import (
	"fmt"
	"log"
	"os"
	"database/sql"

	"github.com/golang-migrate/migrate/v4"
    "github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/joho/godotenv"
	_ "github.com/go-sql-driver/mysql"
    _ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/juancwu/bento/env"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		panic("Could not load .env\n")
	}
	args := os.Args[1:]
	if len(args) < 1 {
		panic("Insufficient argument.")
	}

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

    db, err := sql.Open("mysql", os.Getenv(env.DSN))
    if err != nil {
        panic(err)
    }
    defer db.Close()

    driver, err := mysql.WithInstance(db, &mysql.Config{})
    if err != nil {
        panic(err)
    }

    m, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%s/migrations", dir), "mysql", driver)
	if err != nil {
		panic(err)
	}

	if args[0] == "up" {
		err = m.Up()
	} else if args[0] == "down" {
		err = m.Down()
	}

	if err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	fmt.Printf("Migrations (%s) Done!\n", args[0])
}
