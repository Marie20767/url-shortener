package keys

type KeyValue string

type Key struct {
	Key  KeyValue
	Used bool
}

func (s *KeyStore) GetUnusedKeys() ([]*Key, error) {
	rows, err := s.conn.Query("SELECT key, used FROM keys WHERE used = false")
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

func (s *KeyStore) CreateKey(k KeyValue) error {
	_, err := s.conn.Exec("INSERT INTO keys (key) VALUES ($1)", k)

	if err != nil {
		return err
	}

	return nil
}

func (s *KeyStore) UpdateKey(k KeyValue, u bool) error {
	_, err := s.conn.Exec("UPDATE keys SET used = $1 WHERE = $2", k, u)

	if err != nil {
		return err
	}

	return nil
}
