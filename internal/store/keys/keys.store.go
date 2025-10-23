package keys

import (
	"context"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

type KeyStore struct {
	pool  *pgxpool.Pool
	cache []string
	mu    sync.Mutex
}

func New(ctx context.Context, dbUrl string) (*KeyStore, error) {
	dbPool, err := pgxpool.New(ctx, dbUrl)
	if err != nil {
		return nil, err
	}

	if err := dbPool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	var cache []string

	return &KeyStore{
		pool:  dbPool,
		cache: cache,
	}, nil
}

func (s *KeyStore) Close() {
	s.pool.Close()
}
