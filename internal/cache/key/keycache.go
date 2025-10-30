package keycache

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	client *redis.Client
}

func New(cacheUrl string) (*Cache, error) {
	opt, err := redis.ParseURL(cacheUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to create new key cache: %w", err)
	}

	return &Cache{
		client: redis.NewClient(opt),
	}, nil
}

func (c *Cache) Get(ctx context.Context) (string, bool) {
	key, err := c.client.RandomKey(ctx).Result()
	if err != nil {
		slog.Error("failed to fetch key from cache", slog.Any("error", err))
		return "", false
	}
	if key == "" {
		return "", false
	}

	deleted, err := c.client.Del(ctx, key).Result()
	if err != nil {
		slog.Error("failed to delete used key from cache", slog.Any("error", err))
		return "", false
	}
	if deleted == 0 {
		return "", false
	}

	return key, true
}

func (c *Cache) Add(ctx context.Context, keyMap map[string]string) {
	err := c.client.MSet(ctx, keyMap).Err()
	if err != nil {
		slog.Error("failed to insert keys into cache", slog.Any("error", err))
	}
}
