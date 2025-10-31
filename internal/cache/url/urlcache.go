package urlcache

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
		return nil, fmt.Errorf("failed to create new url cache: %w", err)
	}

	return &Cache{
		client: redis.NewClient(opt),
	}, nil
}

func (c *Cache) Get(ctx context.Context, key string) (string, bool) {
	url, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err != redis.Nil {
			slog.Error("failed to fetch url from cache: ", slog.Any("error", err))
		}
		return "", false
	}

	return url, true
}

func (c *Cache) Add(ctx context.Context, key, value string) {
	// TODO: implement expiry
	err := c.client.Set(ctx, key, value, 0).Err()
	if err != nil {
		slog.Error("failed to insert urls into cache: ", slog.Any("error", err))
	}
}
