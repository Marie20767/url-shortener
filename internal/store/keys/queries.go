package keys

type Key struct {
	Key  string
	Used bool
}

func (s *KeyStore) GetKeys() ([]*Key, error) {
	rows, err := s.conn.Query("SELECT key, used FROM keys WHERE used = true")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keys []*Key

	for rows.Next() {
		k := &Key{}
		if err := rows.Scan(&k.Key, &k.Used); err != nil {
			return nil, err
		}
		keys = append(keys, k)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return keys, nil
}

func (s *KeyStore) CreateKey(key string) error {
	_, err := s.conn.Exec("INSERT INTO keys (key) VALUES ($1)", key)

	if err != nil {
		return err
	}

	return nil
}

func (s *KeyStore) UpdateKey(key string, val bool) error {
	_, err := s.conn.Exec("UPDATE keys SET used = $1 WHERE key = $2", val, key)

	if err != nil {
		return err
	}

	return nil
}