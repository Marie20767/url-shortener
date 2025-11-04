package cron

import (
	"context"
	"log/slog"

	"github.com/robfig/cron/v3"

	"github.com/Marie20767/url-shortener/internal/store/keys"
)

type Cron struct {
	client   *cron.Cron
	keyStore *keys.KeyStore
	schedule string
}

func New(store *keys.KeyStore, schedule string) *Cron {
	return &Cron{
		client:   cron.New(),
		keyStore: store,
		schedule: schedule,
	}
}

func (c *Cron) Add(ctx context.Context) error {
	if err := c.keyStore.GenerateAndStoreKeys(ctx); err != nil {
		return err
	}

	_, err := c.client.AddFunc(c.schedule, func() {
		if keyErr := c.keyStore.GenerateAndStoreKeys(ctx); keyErr != nil {
			slog.Error(keyErr.Error())
		}
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *Cron) Start() {
	c.client.Start()
	slog.Info("key cron scheduler started", slog.String("schedule", c.schedule))
}

func (c *Cron) Stop() context.Context {
	return c.client.Stop()
}
