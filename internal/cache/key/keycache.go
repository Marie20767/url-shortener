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

// lua script needed to ensure atomic key fetching from the cache across all server instances
var getAndDelScript = redis.NewScript(`
	local key = redis.call("RANDOMKEY")
	if not key then
			return nil
	end
	local deleted = redis.call("DEL", key)
	if deleted == 1 then
			return key
	else
			return nil
	end
`)

func (c *Cache) Get(ctx context.Context) (string, bool) {
	res, err := getAndDelScript.Run(ctx, c.client, nil).Result()
	if err != nil {
		slog.Error("failed to fetch key from cache", slog.Any("error", err))
		return "", false
	}
	if res == nil {
		return "", false
	}

	key, ok := res.(string)
	if !ok {
	    slog.Error("unexpected result type from key cache", slog.Any("res", res))
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
