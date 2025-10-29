package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	client *redis.Client
}

func New(cacheUrl string) (*Cache, error) {
	opt, err := redis.ParseURL(cacheUrl)
	if err != nil {
		return nil, err
	}

	return &Cache{
		client: redis.NewClient(opt),
	}, nil
}

func (c *Cache) Get(ctx context.Context, key string) (string, bool) {
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return "", false
	}

	return val, true
}

func (c *Cache) Add(ctx context.Context, key, value string) bool {
	err := c.client.Set(ctx, key, value, 0).Err()

	return err == nil
}
