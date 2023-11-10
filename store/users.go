package store

type User struct {
	Id            int
	Email         string
	ProviderId    int
	ProviderToken *string
}

func (s *Store) CreateNewUser(login string, email string, providerId int) (*User, error) {
	res, err := s.db.Exec("INSERT INTO users (login, email, provider_id) VALUES(?, ?, ?)", login, email, providerId)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	user := User{}

	err = s.db.QueryRow("SELECT id, email, provider_id, provider_token FROM users WHERE id = ?;", id).
		Scan(&user.Id, &user.Email, &user.ProviderId, &user.ProviderToken)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Store) GetUserByGhId(providerId int) (*User, error) {
	var user User
	err := s.db.
		QueryRow("SELECT id, email, provider_id, provider_token FROM users WHERE provider_id = ?", providerId).
		Scan(&user.Id, &user.Email, &user.ProviderId, &user.ProviderToken)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Store) GetUserById(id int) (*User, error) {
	var user User
	err := s.db.QueryRow("SELECT id, email, provider_id, provider_token FROM users WHERE id = ?", id).
		Scan(&user.Id, &user.Email, &user.ProviderId, &user.ProviderToken)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Store) UpdateUserProviderToken(id int, token string) error {
	_, err := s.db.Exec("UPDATE users SET provider_token = ? WHERE id = ?", token, id)
	if err != nil {
		return err
	}

	return nil
}
