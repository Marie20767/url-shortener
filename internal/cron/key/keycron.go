package keycron

import (
	"context"
	"log/slog"

	"github.com/robfig/cron/v3"

	"github.com/Marie20767/url-shortener/internal/keygenerator"
	"github.com/Marie20767/url-shortener/internal/store/keys"
)

type Cron struct {
	client       *cron.Cron
	keyGenerator *keygenerator.KeyGenStore
	schedule     string
}

func New(store *keys.KeyStore, schedule string) *Cron {
	return &Cron{
		client:       cron.New(),
		keyGenerator: keygenerator.New(store),
		schedule:     schedule,
	}
}

func (c *Cron) Add(ctx context.Context) error {
	if err := c.keyGenerator.Run(ctx); err != nil {
		return err
	}

	_, err := c.client.AddFunc(c.schedule, func() {
		if cronErr := c.keyGenerator.Run(ctx); cronErr != nil {
			slog.Error(cronErr.Error())
		}
	})
	if err != nil {
		return err
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
