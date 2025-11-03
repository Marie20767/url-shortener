package keysgenerator

import (
	"context"
	"crypto/rand"
	"fmt"

	"github.com/Marie20767/url-shortener/internal/store/keys"
	"github.com/Marie20767/url-shortener/internal/utils/set"
)

const (
	batchSize         = 50
	keyLength         = 8
	alphanumericChars = "abcdefghijklmnopqrstuvwxyz0123456789"
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
