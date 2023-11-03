package store

import (
    "database/sql"
)

type Store struct {
    db *sql.DB
}

func (s *Store) Stats() sql.DBStats {
    return s.db.Stats()
}
