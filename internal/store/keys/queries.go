package keys

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func (s *KeyStore) GetUnused(ctx context.Context) (string, error) {
	var claimedKey string
	query := `WITH key AS (SELECT key_value FROM keys WHERE used = false LIMIT 1)
						UPDATE keys
						SET used = true
						FROM key
						WHERE keys.key = key.key_value
						RETURNING key.key_value`

	if err := s.pool.QueryRow(ctx, query).Scan(&claimedKey); err != nil {
		return "", err
	}

	return claimedKey, nil
}

func (s *KeyStore) Insert(ctx context.Context, keys []string) (int, error) {
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

func (s *KeyStore) Update(ctx context.Context, used bool, key string) error {
	_, err := s.pool.Exec(ctx, "UPDATE keys SET used = $1 WHERE key_value = $2", used, key)
	if err != nil {
		return err
	}

	return nil
}
