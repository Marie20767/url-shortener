package urlhandlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Marie20767/url-shortener/internal/store/keys"
	"github.com/Marie20767/url-shortener/internal/store/urls"
	"github.com/labstack/echo/v4"
)

type UrlHandler struct {
	UrlStore  *urls.UrlStore
	KeyStore  *keys.KeyStore
	ApiDomain string
}

type UrlData struct {
	Url    string     `json:"url" validate:"required,url"`
	Expiry *time.Time `json:"expiry,omitempty"`
}

type KeyParam struct {
	Key string `param:"key" validate:"required,alphanum,len=8"`
}

func (h *UrlHandler) CreateShort(echoCtx echo.Context) error {
	ctx := echoCtx.Request().Context()

	var req UrlData
	if err := echoCtx.Bind(&req); err != nil {
		return validationErr()
	}

	if err := echoCtx.Validate(&req); err != nil {
		return validationErr()
	}

	tx, err := h.KeyStore.BeginTransaction(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to start transaction")
	}
	defer tx.Rollback(ctx)
	key, err := h.KeyStore.GetUnused(ctx, tx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get unused key")
	}

	urlData := &urls.UrlData{Key: key, Url: req.Url, Expiry: req.Expiry}
	if err := h.UrlStore.Insert(ctx, urlData); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to insert new url data")
	}

	if err := tx.Commit(ctx); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to commit transaction")
	}

	return echoCtx.JSON(http.StatusOK, map[string]string{
		"url": fmt.Sprintf("%s/%s", h.ApiDomain, key),
	})
}

func (h *UrlHandler) GetLong(ctx echo.Context) error {
	var param KeyParam
	if err := ctx.Bind(&param); err != nil {
		return err
	}

	if err := ctx.Validate(&param); err != nil {
		return validationErr()
	}

	longUrl, err := h.UrlStore.Get(ctx.Request().Context(), strings.ToLower(param.Key))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get long url")
	}

	return ctx.Redirect(http.StatusMovedPermanently, longUrl)
}

func validationErr() error {
	return echo.NewHTTPError(http.StatusBadRequest, "Validation Error")
}
