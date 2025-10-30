package keycron

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/robfig/cron/v3"

	"github.com/Marie20767/url-shortener/internal/keygenerator"
	"github.com/Marie20767/url-shortener/internal/store/keys"
)

type Cron struct {
	client   *cron.Cron
	store    *keys.KeyStore
	schedule string
}

func New(store *keys.KeyStore, schedule string) *Cron {
	return &Cron{
		client:   cron.New(),
		store:    store,
		schedule: schedule,
	}
}

func (c *Cron) Add(ctx context.Context) error {
	c.generateKeys(ctx)

	_, err := c.client.AddFunc(c.schedule, func() {
		c.generateKeys(ctx)
	})
	if err != nil {
		return fmt.Errorf("failed to add cron: %w", err)
	}

	return nil
}

func (c *Cron) Start() {
	c.client.Start()
	slog.Info("cron scheduler started", slog.String("schedule", c.schedule))
}

func (c *Cron) Stop() context.Context {
	return c.client.Stop()
}

func (c *Cron) generateKeys(ctx context.Context) {
	keyGen := keygenerator.New(c.store)
	if err := keyGen.Run(ctx); err != nil {
		slog.Error("failed to generate keys", slog.Any("error", err))
	}
}
