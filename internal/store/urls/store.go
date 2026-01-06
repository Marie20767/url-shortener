package urls

import (
	"context"
	"fmt"

	"github.com/Marie20767/url-shortener/internal/store/keys"
	"github.com/Marie20767/url-shortener/internal/utils/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UrlStore struct {
	pool     *pgxpool.Pool
	urlCache *Cache
	keyCache *keys.Cache
}

func New(ctx context.Context, cfg *config.Url, keyCache *keys.Cache) (*UrlStore, error) {
	dbPool, err := pgxpool.New(ctx, cfg.DbUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to create new db pool: %w", err)
	}

	if err := dbPool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	newCache, err := NewCache(cfg.CacheUrl)
	if err != nil {
		return nil, err
	}

	return &UrlStore{
		pool:     dbPool,
		urlCache: newCache,
		keyCache: keyCache,
	}, nil
}

func (s *UrlStore) Ping(ctx context.Context) error {
	return s.pool.Ping(ctx)
}

func (s *UrlStore) Close() {
	s.pool.Close()
}
