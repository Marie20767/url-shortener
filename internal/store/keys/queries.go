package keys

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type KeyValue string

type key struct {
	key  KeyValue
	used bool
}

func (s *KeyStore) GetUnusedKeys(ctx context.Context) ([]*key, error) {
	rows, err := s.pool.Query(ctx, "SELECT key_value, used FROM keys WHERE used = false")
	if err != nil {
		return nil, err
	}

	keys, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (*key, error) {
		k := &key{}
		if err := row.Scan(&k.key, &k.used); err != nil {
			return nil, err
		}
		return k, nil
	})
	if err != nil {
		return nil, err
	}

	return keys, nil
}

func (s *KeyStore) CreateKey(ctx context.Context, key KeyValue) error {
	_, err := s.pool.Exec(ctx, "INSERT INTO keys (key_value) VALUES ($1)", key)
	if err != nil {
		return err
	}

	return nil
}

func (s *KeyStore) UpdateKey(ctx context.Context, used bool, key KeyValue) error {
	_, err := s.pool.Exec(ctx, "UPDATE keys SET used = $1 WHERE key_value = $2", used, key)
	if err != nil {
		return err
	}

	return nil
}
