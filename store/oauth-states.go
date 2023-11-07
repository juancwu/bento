package store

import (
	"strconv"
	"time"
)

func (s *Store) SaveState(stateId string, flow string, redirect string, port string) error {
    exp := time.Now().Add(10 * time.Minute).UTC()
    portInt, err := strconv.Atoi(port)
    if err != nil {
        return err
    }
    portNum := uint16(portInt)
    _, err = s.db.Exec("INSERT INTO oauth_states (state_id, flow, redirect, port, expires_at) VALUES (?, ?, ?, ?, ?);", stateId, flow, redirect, portNum, exp)
    if err != nil {
        return err
    }

    return nil
}
