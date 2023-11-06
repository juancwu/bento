package store

import (
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

type User struct {
    Id int
    Email string
    ObjectId string
}

func (s *Store) CreateNewUser(email string, ghId int) (*User, error) {
    tx, err := s.db.Begin()
    if err != nil {
        return nil, err
    }
    defer tx.Rollback()

    stmt, err := tx.Prepare("INSERT INTO users (email, gh_id, object_id, created_at) VALUES(?, ?, ?, ?)")
    if err != nil {
        return nil, err
    }
    defer stmt.Close()

    nanoid, err := gonanoid.New(12)
    if err != nil {
        return nil, err
    }

    res, err := stmt.Exec(email, ghId, nanoid, time.Now().UTC())
    if err != nil {
        return nil, err
    }

    lastId, err := res.LastInsertId()
    if err != nil {
        return nil, err
    }

    var user User
    err = tx.QueryRow("SELECT id, email, object_id FROM users WHERE id = ?", lastId).Scan(&user.Id, &user.Email, &user.ObjectId)
    if err != nil {
        return nil, err
    }

    err = tx.Commit()
    if err != nil {
        return nil, err
    }

    return &user, nil
}

func (s *Store) CheckUser(ghId int) (*User, error) {
    var user User
    err := s.db.QueryRow("SELECT id, email, object_id FROM users WHERE gh_id = ?", ghId).Scan(&user.Id, &user.Email, &user.ObjectId)
    if err != nil {
        return nil, err
    }

    return &user, nil
}

func (s *Store) GetUserById(id int) (*User, error) {
    var user User
    err := s.db.QueryRow("SELECT id, email, object_id FROM users WHERE id = ?", id).Scan(&user.Id, &user.Email, &user.ObjectId)
    if err != nil {
        return nil, err
    }

    return &user, nil
}
