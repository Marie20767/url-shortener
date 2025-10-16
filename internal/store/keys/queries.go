package keys

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type key struct {
	key  string
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

func (s *KeyStore) InsertKeys(ctx context.Context, keys []string) (int, error) {
	batch := &pgx.Batch{}

	for _, key := range keys {
		batch.Queue("INSERT INTO keys (key_value) VALUES ($1) ON CONFLICT DO NOTHING", key)
	}

	results := s.pool.SendBatch(ctx, batch)
	defer results.Close()

	var totalInserted int64
	for range keys {
		cmdTag, err := results.Exec()
		if err != nil {
			return 0, err
		}
		totalInserted += cmdTag.RowsAffected()
	}

	return int(totalInserted), nil
}

func (s *KeyStore) UpdateKey(ctx context.Context, used bool, key string) error {
	_, err := s.pool.Exec(ctx, "UPDATE keys SET used = $1 WHERE key_value = $2", used, key)
	if err != nil {
		return err
	}

	return nil
}
