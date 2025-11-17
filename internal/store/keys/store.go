package keys

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Marie20767/url-shortener/internal/utils/config"
)

type KeyStore struct {
	pool  *pgxpool.Pool
	cache *Cache
}

func New(ctx context.Context, cfg *config.Key) (*KeyStore, error) {
	dbPool, err := pgxpool.New(ctx, cfg.DbUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to create new key db pool: %w", err)
	}

	if err := dbPool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect to key db: %w", err)
	}

	newCache, err := NewCache(cfg.CacheUrl)
	if err != nil {
		return nil, err
	}

	return &KeyStore{
		pool:  dbPool,
		cache: newCache,
	}, nil
}

func (s *KeyStore) Ping(ctx context.Context) error {
	return s.pool.Ping(ctx)
}

func (s *KeyStore) PingCache(ctx context.Context) error {
	return s.cache.Ping(ctx)
}

func (s *KeyStore) Close() {
	s.pool.Close()
}
