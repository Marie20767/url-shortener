package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

const storeErrStatus = "unreachable"

func (h *Handler) HealthCheck(e echo.Context) error {
	ctx := e.Request().Context()
	keyStoreStatus, urlStoreStatus, cacheStatus := "ok", "ok", "ok"

	err := h.KeyStore.Ping(ctx)
	if err != nil {
		keyStoreStatus = storeErrStatus
	}

	err = h.UrlStore.Ping(ctx)
	if err != nil {
		urlStoreStatus = storeErrStatus
	}

	err = h.KeyStore.PingCache(ctx)
	if err != nil {
		cacheStatus = storeErrStatus
	}

	return e.JSON(http.StatusOK, map[string]string{
		"status": "ok",
		"keydb":  keyStoreStatus,
		"urldb":  urlStoreStatus,
		"cache":  cacheStatus,
	})
}
