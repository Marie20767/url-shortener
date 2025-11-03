package keysgenerator

import (
	"context"
	"crypto/rand"
	"fmt"
	"log/slog"

	"github.com/Marie20767/url-shortener/internal/store/keys"
	"github.com/Marie20767/url-shortener/internal/utils/set"
)

const (
	alphanumericChars = "abcdefghijklmnopqrstuvwxyz0123456789"
	batchSize         = 50
	cacheMinSize      = 10
	keyLength         = 8
)

type KeyGenStore struct {
	keyStore *keys.KeyStore
}

func New(keyStore *keys.KeyStore) *KeyGenStore {
	return &KeyGenStore{
		keyStore: keyStore,
	}
}

func (s *KeyGenStore) Run(ctx context.Context) error {
	currentCacheSize := s.keyStore.GetCacheSize(ctx)
	if currentCacheSize >= cacheMinSize {
		slog.Info("cache", "size", currentCacheSize)
		slog.Info("enough keys in cache, skipped key generation")
		return nil
	}

	rowsInserted := 0
	for rowsInserted < batchSize {
		newKeys := make([]string, 0, batchSize)

		for range batchSize {
			random, err := randomString(keyLength)
			if err != nil {
				return fmt.Errorf("failed to generate keys: %w", err)
			}
			newKeys = append(newKeys, random)
		}

		keysWithoutDuplicates := set.New(newKeys...).ToSlice()

		rows, err := s.keyStore.Insert(ctx, keysWithoutDuplicates)
		if err != nil {
			return fmt.Errorf("failed to insert newly generated keys: %w", err)
		}

		rowsInserted += rows
	}
	slog.Info("successfully generated keys!")

	return nil
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
