package urls

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/Marie20767/url-shortener/internal/store/urls/model"
)

type Cache struct {
	client *redis.Client
}

func NewCache(cacheUrl string) (*Cache, error) {
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

func (c *Cache) Add(ctx context.Context, urlData *model.UrlData, currentTimestamp time.Time) {
	var expiry time.Duration
	switch urlData.Expiry {
	case nil:
		expiry = 0
	default:
		expiry = currentTimestamp.Sub(*urlData.Expiry)
	}

	err := c.client.Set(ctx, urlData.Key, urlData.Url, expiry).Err()
	if err != nil {
		slog.Error("failed to insert urls into cache: ", slog.Any("error", err))
	}
}
