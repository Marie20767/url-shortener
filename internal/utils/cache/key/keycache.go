package keycache

import (
	"context"
	"fmt"

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

func (c *Cache) Get(ctx context.Context) (string, bool) {
	key, err := c.client.RandomKey(ctx).Result()
	if err != nil {
		fmt.Println(">>> failed to fetch key from cache: ", err)
		return "", false
	}
	if key == "" {
		return "", false
	}

	deleted, err := c.client.Del(ctx, key).Result()
	if err != nil {
		fmt.Println(">>> failed to delete used key from cache: ", err)
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
		fmt.Println(">>> failed to insert keys into cache: ", err)
	}
}
