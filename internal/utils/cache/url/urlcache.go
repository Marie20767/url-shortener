package urlcache

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

func (c *Cache) Get(ctx context.Context, key string) (string, bool) {
	url, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err != redis.Nil {
			fmt.Println(">>> failed to fetch long url from cache: ", err)
		}
		return "", false
	}

	return url, true
}

func (c *Cache) Add(ctx context.Context, key, value string) {
	// TODO: implement expiry
	err := c.client.Set(ctx, key, value, 0).Err()
	if err != nil {
		fmt.Println(">>> failed to insert keys into cache: ", err)
	}
}
