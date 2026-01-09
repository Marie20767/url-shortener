package jobs

import (
	"context"
	"log/slog"

	"github.com/Marie20767/url-shortener/internal/store/urls"
)

func KeyGenerationJob(urlStore *urls.UrlStore) func(ctx context.Context) {
	return func(ctx context.Context) {
		if err := urlStore.GenerateAndStoreKeys(ctx); err != nil {
			slog.Error(err.Error())
		}
	}
}

func UrlCleanUpJob(urlStore *urls.UrlStore) func(ctx context.Context) {
	return func(ctx context.Context) {
		deletedKeys := cleanupUrls(urlStore, ctx)
		freeUpKeys(urlStore, ctx, deletedKeys)
	}
}

func cleanupUrls(urlStore *urls.UrlStore, ctx context.Context) []string {
	deletedKeys, err := urlStore.DeleteExpiredUrls(ctx)

	switch {
	case err != nil:
		slog.Error("url-cleanup: failed to delete urls", slog.Any("error", err))
	case len(deletedKeys) > 0:
		slog.Debug("url-cleanup: successfully deleted urls", slog.Int("number", len(deletedKeys)))
	default:
		slog.Debug("url-cleanup: no expired urls to delete")
	}

	return deletedKeys
}

func freeUpKeys(urlStore *urls.UrlStore, ctx context.Context, deletedKeys []string) {
	freedUpKeyCount, err := urlStore.FreeUpUnusedKeys(ctx, deletedKeys)

	switch {
	case err != nil:
		slog.Error("url-cleanup: failed to free up unused keys", slog.Any("error", err))
	case freedUpKeyCount > 0:
		slog.Debug("url-cleanup: successfully freed up unused keys", slog.Int("number", freedUpKeyCount))
	default:
		slog.Debug("url-cleanup: no keys to free up")
	}
}
