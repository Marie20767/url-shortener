package keys

type KeyValue string

type key struct {
	key  KeyValue
	used bool
}

func (s *KeyStore) GetUnusedKeys() ([]*key, error) {
	rows, err := s.conn.Query("SELECT key_value, used FROM keys WHERE used = false")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keys []*key

	for rows.Next() {
		k := &key{}
		if err := rows.Scan(&k.key, &k.used); err != nil {
			return nil, err
		}
		keys = append(keys, k)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return keys, nil
}

func (s *KeyStore) CreateKey(key KeyValue) error {
	_, err := s.conn.Exec("INSERT INTO keys (key_value) VALUES ($1)", key)

	if err != nil {
		return err
	}

	return nil
}

func (s *KeyStore) UpdateKey(used bool, key KeyValue) error {
	_, err := s.conn.Exec("UPDATE keys SET used = $1 WHERE key_value = $2", used, key)

	if err != nil {
		return err
	}

	return nil
}
