package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	keycron "github.com/Marie20767/url-shortener/internal/cron/keys"
	"github.com/Marie20767/url-shortener/internal/cron/model"
	"github.com/Marie20767/url-shortener/internal/server"
	"github.com/Marie20767/url-shortener/internal/store/keys"
	"github.com/Marie20767/url-shortener/internal/store/urls"
	"github.com/Marie20767/url-shortener/internal/utils/config"
)

const serverTimeout = 10

func main() {
	if err := run(); err != nil {
		slog.Error("run failed", slog.Any("error", err))
		os.Exit(1)
	}

	slog.Info("shutting down gracefully...")
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.ParseEnv()
	if err != nil {
		return err
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: cfg.LogLevel,
	}))
	slog.SetDefault(logger)

	keyStore, err := keys.New(ctx, cfg.Key)
	if err != nil {
		return err
	}
	defer keyStore.Close()
	slog.Info("successfully connected to key db!")

	urlStore, err := urls.New(cfg.Url)
	if err != nil {
		return err
	}
	defer urlStore.Close(ctx) //nolint:errcheck
	slog.Info("successfully connected to url db!")

	keyCron := keycron.New(keyStore, cfg.Key.CronSchedule)
	cancelKeyCron, err := setupCron(keyCron)
	defer cancelKeyCron()
	if err != nil {
		return err
	}

	urlCron := urlCron.New(urlStore, cfg.Key.CronSchedule)
	cancelUrlCron, err := setUpCron(urlCron)
	defer cancelUrlCron()
	if err != nil {
		return err
	}

	s := server.New(keyStore, urlStore, cfg.Domain)
	go func() {
		if err := s.Start(cfg.Port); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", slog.Any("error", err))
		}
	}()

	// block until shutdown signal
	<-ctx.Done()
	slog.Info("shutdown signal received")

	// cancel cron contexts to prevent new jobs from starting
	cancelKeyCron()
	cancelUrlCron()

	stopKeyCtx := keyCron.Stop()
	<-stopKeyCtx.Done()
	slog.Info("key cron jobs completed")

	stopUrlCtx := urlCron.Stop()
	<-stopUrlCtx.Done()
	slog.Info("url cron jobs completed")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), serverTimeout*time.Second)
	defer cancel()

	if err := s.Stop(shutdownCtx); err != nil {
		slog.Error("server shutdown error", slog.Any("error", err))
	}

	return nil
}

func setupCron(cron model.CronLike) (context.CancelFunc, error) {
	cronCtx, cancelCron := context.WithCancel(context.Background())

	err := cron.Add(cronCtx)
	if err != nil {
		return cancelCron, err
	}
	cron.Start()

	return cancelCron, nil
}
