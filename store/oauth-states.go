package store

import (
	"time"
)

func (s *Store) SaveState(stateId string, flow string, redirect string, port uint16) error {
    exp := time.Now().Add(10 * time.Minute).UTC()
    _, err := s.db.Exec("INSERT INTO oauth_states (state_id, flow, redirect, port, expires_at) VALUES (?, ?, ?, ?, ?);", stateId, flow, redirect, port, exp)
    if err != nil {
        return err
    }

    return nil
}
