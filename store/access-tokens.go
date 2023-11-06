package store

func (s *Store) GetAccessToken(id int) (string, error) {
    var token string
    err := s.db.QueryRow("SELECT token FROM gh_access_tokens WHERE user_id = ?", id).Scan(&token)
    if err != nil {
        return "", err
    }

    return token, nil
}

func (s *Store) UpdateAccessToken(id int, token string) error {
    _, err := s.db.Exec("UPDATE gh_access_tokens SET token = ? WHERE user_id = ?", token, id)
    return err
}

func (s *Store) SaveAccessToken(id int, token string) error {
    _, err := s.db.Exec("INSERT INTO gh_access_tokens (token, user_id) VALUES (?, ?)", token, id)
    return err
}
