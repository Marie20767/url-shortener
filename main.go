package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	urlhandlers "github.com/Marie20767/url-shortener/api/handlers/url"
	"github.com/Marie20767/url-shortener/api/routes"
	"github.com/Marie20767/url-shortener/internal/cron/key"
	"github.com/Marie20767/url-shortener/internal/store/keys"
	"github.com/Marie20767/url-shortener/internal/store/urls"
	"github.com/Marie20767/url-shortener/internal/utils/config"
)

const serverTimeout = 10

type stores struct {
	keyStore *keys.KeyStore
	urlStore *urls.UrlStore
}

type customValidator struct {
	validator *validator.Validate
}

func (cv *customValidator) Validate(i any) error {
	return cv.validator.Struct(i)
}

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
	cancelCron, err := setupKeyCron(keyCron)
	defer cancelCron()
	if err != nil {
		return err
	}

	e := setupServer(&stores{
		keyStore: keyStore,
		urlStore: urlStore,
	}, cfg.Domain)

	go func() {
		if err := e.Start(":" + cfg.Port); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", slog.Any("error", err))
		}
	}()

	// block until shutdown signal
	<-ctx.Done()
	slog.Info("shutdown signal received")

	// cancel cron context to prevent new jobs from starting
	cancelCron()

	stopCtx := keyCron.Stop()
	<-stopCtx.Done()
	slog.Info("cron jobs completed")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), serverTimeout*time.Second)
	defer cancel()

	if err := e.Shutdown(shutdownCtx); err != nil {
		slog.Error("server shutdown error", slog.Any("error", err))
	}

	return nil
}

func setupKeyCron(cron *keycron.Cron) (context.CancelFunc, error) {
	kCronCtx, cancelCron := context.WithCancel(context.Background())

	kCronErr := cron.Add(kCronCtx)
	if kCronErr != nil {
		return cancelCron, kCronErr
	}
	cron.Start()

	return cancelCron, nil
}

func setupServer(stores *stores, domain string) *echo.Echo {
	e := echo.New()
	e.Validator = &customValidator{validator: validator.New()}
	urlHandler := &urlhandlers.UrlHandler{
		KeyStore:  stores.keyStore,
		UrlStore:  stores.urlStore,
		ApiDomain: domain,
	}
	routes.RegisterAll(e, urlHandler)

	return e
}
