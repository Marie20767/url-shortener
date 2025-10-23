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

func (h *UrlHandler) CreateShort(ctx echo.Context) error {
	var req UrlData
	if err := ctx.Bind(&req); err != nil {
		return validationErr()
	}

	if err := ctx.Validate(&req); err != nil {
		return validationErr()
	}

	key, keyErr := h.KeyStore.GetUnused(ctx.Request().Context())
	if keyErr != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get unused key")
	}

	urlData := &urls.UrlData{Key: key, Url: req.Url, Expiry: req.Expiry}
	if urlErr := h.UrlStore.Insert(ctx.Request().Context(), urlData); urlErr != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to insert new url data")
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"url": fmt.Sprintf("%s/%s", h.ApiDomain, key),
	})
}

func (h *UrlHandler) GetLong(ctx echo.Context) error {
	key := strings.ToLower(ctx.Param("key"))
	param := KeyParam{Key: key}

	if err := ctx.Validate(&param); err != nil {
		return validationErr()
	}

	longUrl, err := h.UrlStore.Get(ctx.Request().Context(), key)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get long url")
	}

	return ctx.Redirect(http.StatusMovedPermanently, longUrl)
}

func validationErr() error {
	return echo.NewHTTPError(http.StatusBadRequest, "Validation Error")
}
