package keys

import (
	"context"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"

	keycache "github.com/Marie20767/url-shortener/internal/cache/key"
	"github.com/Marie20767/url-shortener/internal/utils/config"
)

type KeyStore struct {
	pool  *pgxpool.Pool
	cache *keycache.Cache
	mu    sync.Mutex
}

func New(ctx context.Context, cfg *config.Key) (*KeyStore, error) {
	dbPool, err := pgxpool.New(ctx, cfg.DbUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to create new key db pool: %w", err)
	}

	if err := dbPool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect to key db: %w", err)
	}

	newCache, err := keycache.New(cfg.CacheUrl)
	if err != nil {
		return nil, err
	}

	return &KeyStore{
		pool:  dbPool,
		cache: newCache,
	}, nil
}

func (s *KeyStore) Close() {
	s.pool.Close()
}
