package keys

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func (s *KeyStore) BeginTransaction(ctx context.Context) (pgx.Tx, error) {
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})

	return tx, err
}

func (s *KeyStore) GetUnused(ctx context.Context, tx pgx.Tx) (string, error) {
	var claimedKey string

	if len(s.cache) > 0 {
    s.mu.Lock()
		claimedKey, s.cache = s.cache[0], s.cache[1:]
		s.mu.Unlock()

		if err := tx.QueryRow(ctx, "UPDATE keys SET used = true WHERE keys.key_value = $1", claimedKey); err != nil {
			return "", nil
		}

		return claimedKey, nil
	}

	query := `WITH key AS (SELECT key_value FROM keys WHERE used = false LIMIT 1)
						UPDATE keys
						SET used = true
						FROM key
						WHERE keys.key_value = key.key_value
						RETURNING key.key_value`

	if err := tx.QueryRow(ctx, query).Scan(&claimedKey); err != nil {
		return "", err
	}

	return claimedKey, nil
}

func (s *KeyStore) Insert(ctx context.Context, keys []string) (int, error) {
	query := "INSERT INTO keys (key_value) SELECT UNNEST($1::varchar(8)[]) ON CONFLICT DO NOTHING RETURNING key_value"
	rows, err := s.pool.Query(ctx, query, keys)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var inserted []string
	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err != nil {
			return 0, err
		}

		inserted = append(inserted, key)
	}

	s.cache = append(s.cache, inserted...)

	return len(inserted), nil
}

func (s *KeyStore) Update(ctx context.Context, used bool, key string) error {
	_, err := s.pool.Exec(ctx, "UPDATE keys SET used = $1 WHERE key_value = $2", used, key)
	if err != nil {
		return err
	}

	return nil
}
