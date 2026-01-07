package urls

import (
	"context"
	"crypto/rand"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/Marie20767/url-shortener/internal/store/urls/model"
	"github.com/Marie20767/url-shortener/internal/utils/set"
)

const (
	alphanumericChars = "aAbBcCdDeEfFgGhHiIjJkKlLmMnNoOpPqQrRsStTuUvVwWxXyYzZ0123456789"
	keyBatchSize      = 50
	keyLength         = 8
)

func (s *UrlStore) GetUnusedKey(ctx context.Context, tx pgx.Tx) (string, error) {
	if s.keyCache.ShouldRefill(ctx) {
		// prevents blocking response while new keys are being generated
		go func() {
			if err := s.GenerateAndStoreKeys(ctx); err != nil {
				// in future this could trigger an alert
				slog.Error(err.Error())
			}
		}()
	}

	claimedKey, ok := s.keyCache.Get(ctx)
	if ok {
		_, err := tx.Exec(ctx, "UPDATE keys SET used = true WHERE id = $1", claimedKey)
		if err != nil {
			return "", fmt.Errorf("failed to update used key in db: %w", err)
		}

		return claimedKey, nil
	}

	query := `WITH key AS (
    SELECT id 
    FROM keys 
    WHERE used = false 
    FOR UPDATE SKIP LOCKED 
    LIMIT 1
	)
	UPDATE keys
	SET used = true
	FROM key
	WHERE keys.id = key.id
	RETURNING key.id`

	if err := tx.QueryRow(ctx, query).Scan(&claimedKey); err != nil {
		return "", fmt.Errorf("failed to fetch unused key from db: %w", err)
	}

	return claimedKey, nil
}

func (s *UrlStore) FreeUpUnusedKeys(ctx context.Context, keys []string) (int, error) {
	batch := &pgx.Batch{}
	for _, key := range keys {
		batch.Queue("UPDATE keys SET used = false WHERE id = $1", key)
	}
	results := s.pool.SendBatch(ctx, batch)
	defer results.Close() //nolint:errcheck

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

func (s *UrlStore) InsertNewKeys(ctx context.Context, keys []string) (int, error) {
	// do nothing on conflict because we just want to ignore any duplicate keys
	query := "INSERT INTO keys (id) SELECT UNNEST($1::varchar(8)[]) ON CONFLICT DO NOTHING RETURNING id"
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
	s.keyCache.Set(ctx, inserted)

	return len(inserted), nil
}

func (s *UrlStore) GenerateAndStoreKeys(ctx context.Context) error {
	if !s.keyCache.ShouldRefill(ctx) {
		return nil
	}

	rowsInserted := 0
	for rowsInserted < keyBatchSize {
		newKeys, err := generateKeys()
		if err != nil {
			return err
		}

		rows, err := s.InsertNewKeys(ctx, newKeys)
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
		key, err := getRandomString(keyLength)
		if err != nil {
			return nil, fmt.Errorf("failed to generate keys: %w", err)
		}
		newKeys = append(newKeys, key)
	}

	return set.New(newKeys...).ToSlice(), nil
}

func getRandomString(length int) (string, error) {
	bytes := make([]byte, length)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	for i, b := range bytes {
		bytes[i] = alphanumericChars[int(b)%len(alphanumericChars)]
	}

	return string(bytes), nil
}

func (s *UrlStore) InsertNewUrl(ctx context.Context, tx pgx.Tx, urlData *model.UrlData) error {
	if urlData.Expiry != nil {
		utcTime := urlData.Expiry.UTC()
		urlData.Expiry = &utcTime
	}

	query := "INSERT INTO urls (short, long, expiry) VALUES ($1, $2, $3)"
	_, err := tx.Exec(ctx, query, urlData.Key, urlData.Url, urlData.Expiry)
	if err != nil {
		return fmt.Errorf("failed to insert new url into db: %w", err)
	}

	return nil
}

func (s *UrlStore) DeleteExpiredUrls(ctx context.Context) ([]string, error) {
	var deletedKeys []string
	query := "DELETE FROM urls WHERE expiry <= NOW() RETURNING short"
	rows, err := s.pool.Query(ctx, query)
	if err != nil {
		return deletedKeys, fmt.Errorf("failed to delete expired urls: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err != nil {
			return deletedKeys, err
		}

		deletedKeys = append(deletedKeys, key)
	}

	return deletedKeys, nil
}

func (s *UrlStore) GetLongUrl(ctx context.Context, key string) (string, error) {
	now := time.Now().UTC()
	url, ok := s.urlCache.Get(ctx, key)
	if ok {
		return url, nil
	}

	var urlData model.UrlData
	query := "SELECT short, long, expiry FROM urls WHERE short = $1 AND (expiry > NOW() OR expiry IS NULL)"
	row := s.pool.QueryRow(ctx, query, key)
	if err := row.Scan(&urlData.Key, &urlData.Url, &urlData.Expiry); err != nil {
		return "", err
	}

	s.urlCache.Set(ctx, &urlData, now)

	return urlData.Url, nil
}

func (s *UrlStore) BeginTransaction(ctx context.Context) (pgx.Tx, error) {
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})

	return tx, err
}
