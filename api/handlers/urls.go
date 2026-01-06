package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/Marie20767/url-shortener/internal/store/urls"
	"github.com/Marie20767/url-shortener/internal/store/urls/model"
)

type UrlData struct {
	Url    string     `json:"url" validate:"required,url"`
	Expiry *time.Time `json:"expiry,omitempty"`
}

type KeyParam struct {
	Key string `param:"key" validate:"required,alphanum,len=8"`
}

func (h *Handler) CreateShort(e echo.Context) error {
	ctx := e.Request().Context()

	var req UrlData
	if err := e.Bind(&req); err != nil {
		return validationErr()
	}

	if err := e.Validate(&req); err != nil {
		return validationErr()
	}

	tx, err := h.UrlStore.BeginTransaction(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to start transaction")
	}
	defer tx.Rollback(ctx) //nolint:errcheck
	key, err := h.UrlStore.GetUnusedKey(ctx, tx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get unused key")
	}

	urlData := &model.UrlData{Key: key, Url: req.Url, Expiry: req.Expiry}
	err = h.UrlStore.InsertNewUrl(ctx, urlData)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to insert new url data")
	}

	if err := tx.Commit(ctx); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to commit transaction")
	}

	return e.JSON(http.StatusCreated, map[string]string{
		"url": fmt.Sprintf("%s/%s", h.ApiDomain, key),
	})
}

func (h *Handler) GetLong(e echo.Context) error {
	var param KeyParam
	if err := e.Bind(&param); err != nil {
		return err
	}

	if err := e.Validate(&param); err != nil {
		return validationErr()
	}

	longUrl, err := h.UrlStore.GetLongUrl(e.Request().Context(), strings.ToLower(param.Key))
	if err != nil {
		if err == urls.ErrNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "url not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get url")
	}

	return e.Redirect(http.StatusFound, longUrl)
}

func validationErr() error {
	return echo.NewHTTPError(http.StatusBadRequest, "validation Error")
}
