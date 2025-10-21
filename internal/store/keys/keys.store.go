package keys

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type KeyStore struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context, dbUrl string) (*KeyStore, error) {
	dbPool, err := pgxpool.New(ctx, dbUrl)
	if err != nil {
		return nil, err
	}

	if err := dbPool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	return &KeyStore{
		pool: dbPool,
	}, nil
}

func (s *KeyStore) Close() {
	s.pool.Close()
}
