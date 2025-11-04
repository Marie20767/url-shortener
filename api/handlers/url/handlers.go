package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/Marie20767/url-shortener/internal/store/keys"
	"github.com/Marie20767/url-shortener/internal/store/urls"
	"github.com/Marie20767/url-shortener/internal/store/urls/model"
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
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to start transaction")
	}
	defer tx.Rollback(ctx) //nolint:errcheck
	key, err := h.KeyStore.GetUnused(ctx, tx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get unused key")
	}

	urlData := &model.UrlData{Key: key, Url: req.Url, Expiry: req.Expiry}
	id, err := h.UrlStore.Insert(ctx, urlData)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to insert new url data")
	}

	if err := tx.Commit(ctx); err != nil {
		_ = h.UrlStore.DeleteById(ctx, id)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to commit transaction")
	}

	return echoCtx.JSON(http.StatusOK, map[string]string{
		"url": fmt.Sprintf("%s/%s", h.ApiDomain, key),
	})
}

func (h *UrlHandler) GetLong(ctx echo.Context) error {
	now := time.Now().UTC()
	var param KeyParam
	if err := ctx.Bind(&param); err != nil {
		return err
	}

	if err := ctx.Validate(&param); err != nil {
		return validationErr()
	}

	longUrl, err := h.UrlStore.Get(ctx.Request().Context(), strings.ToLower(param.Key), now)
	if err != nil {
		if err == urls.ErrNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "url not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get url")
	}

	return ctx.Redirect(http.StatusMovedPermanently, longUrl)
}

func validationErr() error {
	return echo.NewHTTPError(http.StatusBadRequest, "validation Error")
}
