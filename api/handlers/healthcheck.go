package handlers

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

const storeErrStatus = "unreachable"

func (h *Handler) HealthCheck(e echo.Context) error {
	ctx := e.Request().Context()

	check := func(ping func(context.Context) error) string {
		if err := ping(ctx); err != nil {
			return storeErrStatus
		}
		return "ok"
	}

	keyStoreStatus := check(h.KeyStore.Ping)
	urlStoreStatus := check(h.UrlStore.Ping)
	cacheStatus := check(h.KeyStore.PingCache)

	httpStatus := http.StatusOK
	if keyStoreStatus != "ok" || urlStoreStatus != "ok" || cacheStatus != "ok" {
		httpStatus = http.StatusServiceUnavailable
	}

	return e.JSON(httpStatus, map[string]string{
		"status": "ok",
		"key_db": keyStoreStatus,
		"url_db": urlStoreStatus,
		"cache":  cacheStatus,
	})
}
