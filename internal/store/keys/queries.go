package keys

import (
	"context"
	"crypto/rand"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"

	"github.com/Marie20767/url-shortener/internal/utils/set"
)

const (
	alphanumericChars = "abcdefghijklmnopqrstuvwxyz0123456789"
	keyBatchSize      = 50
	keyLength         = 8
)

func (s *KeyStore) BeginTransaction(ctx context.Context) (pgx.Tx, error) {
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})

	return tx, err
}

func (s *KeyStore) GetUnused(ctx context.Context, tx pgx.Tx) (string, error) {
	if s.cache.ShouldRefillCache(ctx) {
		// prevents blocking response while new keys are being generated
		go func() {
			if err := s.GenerateAndStoreKeys(ctx); err != nil {
				// in future this could trigger an alert
				slog.Error(err.Error())
			}
		}()
	}

	claimedKey, ok := s.cache.Get(ctx)
	if ok {
		_, err := tx.Exec(ctx, "UPDATE keys SET used = true WHERE key_value = $1", claimedKey)
		if err != nil {
			return "", fmt.Errorf("failed to update used key in db: %w", err)
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
		return "", fmt.Errorf("failed to fetch & update used key from db: %w", err)
	}

	return claimedKey, nil
}

func (s *KeyStore) FreeUpUnusedKeys(ctx context.Context, keys []string) (int, error) {
	batch := &pgx.Batch{}
	for _, key := range keys {
		batch.Queue("UPDATE keys SET used = false WHERE key_value = $1", key)
	}
	results := s.pool.SendBatch(ctx, batch)
	defer results.Close()

	count := 0
	for range keys {
		_, err := results.Exec()
		if err != nil {
			return count, err
		}
		count++
	}

	return count, nil
}

func (s *KeyStore) Insert(ctx context.Context, keys []string) (int, error) {
	// do nothing on conflict because we just want to ignore any duplicate keys
	query := "INSERT INTO keys (key_value) SELECT UNNEST($1::varchar(8)[]) ON CONFLICT DO NOTHING RETURNING key_value"
	rows, err := s.pool.Query(ctx, query, keys)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	inserted := make(map[string]string)
	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err != nil {
			return 0, err
		}

		inserted[key] = key
	}
	s.cache.Add(ctx, inserted)

	return len(inserted), nil
}

func (s *KeyStore) GenerateAndStoreKeys(ctx context.Context) error {
	if !s.cache.ShouldRefillCache(ctx) {
		return nil
	}

	rowsInserted := 0
	for rowsInserted < keyBatchSize {
		newKeys, err := generateKeys()
		if err != nil {
			return err
		}

		rows, err := s.Insert(ctx, newKeys)
		if err != nil {
			return fmt.Errorf("failed to insert newly generated keys: %w", err)
		}

		rowsInserted += rows
	}

	slog.Debug("successfully generated keys!")
	return nil
}

func generateKeys() ([]string, error) {
	newKeys := make([]string, 0, keyBatchSize)
	for range keyBatchSize {
		key, err := randomString(keyLength)
		if err != nil {
			return nil, fmt.Errorf("failed to generate keys: %w", err)
		}
		newKeys = append(newKeys, key)
	}

	return set.New(newKeys...).ToSlice(), nil
}

func randomString(length int) (string, error) {
	bytes := make([]byte, length)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	for i, b := range bytes {
		bytes[i] = alphanumericChars[int(b)%len(alphanumericChars)]
	}

	return string(bytes), nil
}
