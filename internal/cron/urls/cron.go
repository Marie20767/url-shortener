package cron

import (
	"context"
	"log/slog"

	"github.com/robfig/cron/v3"

	"github.com/Marie20767/url-shortener/internal/store/keys"
	"github.com/Marie20767/url-shortener/internal/store/urls"
)

type Cron struct {
	client   *cron.Cron
	keyStore *keys.KeyStore
	schedule string
	urlStore *urls.UrlStore
}

func New(keyStore *keys.KeyStore, urlStore *urls.UrlStore, schedule string) *Cron {
	return &Cron{
		client:   cron.New(),
		keyStore: keyStore,
		schedule: schedule,
		urlStore: urlStore,
	}
}

func (c *Cron) Add(ctx context.Context) error {
	_, err := c.client.AddFunc(c.schedule, func() {
		c.cleanupExpiredUrls(ctx)
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *Cron) cleanupExpiredUrls(ctx context.Context) {
	deletedKeys, urlErr := c.urlStore.DeleteExpired(ctx)
	switch {
	case urlErr != nil:
		slog.Error("failed to delete all expired urls from db", slog.Any("error", urlErr))
	case len(deletedKeys) > 0:
		slog.Debug("successfully deleted all expired urls", slog.Int("number", len(deletedKeys)))
	default:
		slog.Debug("no expired urls to delete")
	}

	freedUpKeyCount, keyErr := c.keyStore.FreeUpUnusedKeys(ctx, deletedKeys)
	switch {
	case keyErr != nil:
		slog.Error("failed to free up all unused keys in db", slog.Any("error", urlErr))
	case freedUpKeyCount > 0:
		slog.Debug("successfully freed up all unused keys", slog.Int("number", freedUpKeyCount))
	default:
		slog.Debug("no keys to free up")
	}
}

func (c *Cron) Start() {
	c.client.Start()
	slog.Info("url cron scheduler started", slog.String("schedule", c.schedule))
}

func (c *Cron) Stop() context.Context {
	return c.client.Stop()
}
