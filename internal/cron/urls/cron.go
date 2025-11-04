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

func New(urlStore *urls.UrlStore, keyStore *keys.KeyStore, schedule string) *Cron {
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
	if urlErr != nil {
		// in future this could trigger an alert
		slog.Error("failed to delete all expired urls from db", slog.Any("error", urlErr))
	} else {
		slog.Debug("successfully deleted all expired urls", slog.Int("number", len(deletedKeys)))
	}

	freedUpKeyCount, keyErr := c.keyStore.FreeUpUnusedKeys(ctx, deletedKeys)
	if keyErr != nil {
		slog.Error("failed to free up all unused keys in db", slog.Any("error", urlErr))
	} else {
		slog.Debug("successfully freed up all unused keys", slog.Int("number", freedUpKeyCount))
	}
}

func (c *Cron) Start() {
	c.client.Start()
	slog.Info("url cron scheduler started", slog.String("schedule", c.schedule))
}

func (c *Cron) Stop() context.Context {
	return c.client.Stop()
}
