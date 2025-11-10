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
			slog.Error("url-cleanup: failed to delete urls", slog.Any("error", urlErr))
		case len(deletedKeys) > 0:
			slog.Debug("url-cleanup: successfully deleted urls", slog.Int("number", len(deletedKeys)))
		default:
			slog.Debug("url-cleanup: no expired urls to delete")
		}

		freedUpKeyCount, keyErr := keyStore.FreeUpUnusedKeys(ctx, deletedKeys)
		switch {
		case keyErr != nil:
			slog.Error("url-cleanup: failed to free up unused keys", slog.Any("error", urlErr))
		case freedUpKeyCount > 0:
			slog.Debug("url-cleanup: successfully freed up unused keys", slog.Int("number", freedUpKeyCount))
		default:
			slog.Debug("url-cleanup: no keys to free up")
		}
	}
}
