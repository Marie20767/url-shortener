package jobs

import (
	"context"
	"log/slog"

	"github.com/Marie20767/url-shortener/internal/store/keys"
	"github.com/Marie20767/url-shortener/internal/store/urls"
)

func KeyGenerationJob(keyStore *keys.KeyStore) func(ctx context.Context) {
	return func(ctx context.Context) {
		if err := keyStore.GenerateAndStoreKeys(ctx); err != nil {
			slog.Error(err.Error())
		}
	}
}

func UrlCleanUpJob(keyStore *keys.KeyStore, urlStore *urls.UrlStore) func(ctx context.Context) {
	return func(ctx context.Context) {
		deletedKeys, urlErr := urlStore.DeleteExpired(ctx)
		switch {
		case urlErr != nil:
			slog.Error("failed to delete all expired urls from db", slog.Any("error", urlErr))
		case len(deletedKeys) > 0:
			slog.Debug("successfully deleted all expired urls", slog.Int("number", len(deletedKeys)))
		default:
			slog.Debug("no expired urls to delete")
		}

		freedUpKeyCount, keyErr := keyStore.FreeUpUnusedKeys(ctx, deletedKeys)
		switch {
		case keyErr != nil:
			slog.Error("failed to free up all unused keys in db", slog.Any("error", urlErr))
		case freedUpKeyCount > 0:
			slog.Debug("successfully freed up all unused keys", slog.Int("number", freedUpKeyCount))
		default:
			slog.Debug("no keys to free up")
		}
	}
}
