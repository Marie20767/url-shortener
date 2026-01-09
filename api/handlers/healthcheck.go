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

	keyCacheStatus := check(h.KeyCache.Ping)
	urlStoreStatus := check(h.UrlStore.Ping)

	httpStatus := http.StatusOK
	if keyCacheStatus != "ok" || urlStoreStatus != "ok" {
		httpStatus = http.StatusServiceUnavailable
	}

	return e.JSON(httpStatus, map[string]string{
		"status": "ok",
		"cache":  keyCacheStatus,
		"url_db": urlStoreStatus,
	})
}
