package keygenerator

import (
	"context"
	"crypto/rand"
	"errors"

	"github.com/Marie20767/url-shortener/internal/store/keys"
	"github.com/Marie20767/url-shortener/internal/utils/set"
)

const (
	batchSize         = 50
	keyLength         = 8
	alphanumericChars = "abcdefghijklmnopqrstuvwxyz0123456789"
)

type keygenstore struct {
	keyStore *keys.KeyStore
}

func New(keyStore *keys.KeyStore) *keygenstore {
	return &keygenstore{
		keyStore: keyStore,
	}
}

func (s *keygenstore) Generate(ctx context.Context) error {
	rowsInserted := 0

	for rowsInserted < batchSize {
		keys := make([]string, 0, batchSize)

		for range batchSize {
			random, err := randomString(keyLength)
			if err != nil {
				return errors.New("failed to create url keys")
			}
			keys = append(keys, random)
		}

		keysWithoutDuplicates := set.New(keys...).ToSlice()

		rows, err := s.keyStore.Insert(ctx, keysWithoutDuplicates)
		if err != nil {
			break
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
