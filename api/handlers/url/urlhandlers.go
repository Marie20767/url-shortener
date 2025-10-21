package urlhandlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Marie20767/url-shortener/internal/store/keys"
	"github.com/Marie20767/url-shortener/internal/store/urls"
	"github.com/labstack/echo/v4"
)

type UrlHandler struct {
	UrlDb     *urls.UrlStore
	KeyDb     *keys.KeyStore
	ApiDomain string
}

type CreateShortUrlRequest struct {
	Url    string     `json:"url" validate:"required,url"`
	Expiry *time.Time `json:"expiry,omitempty" validate:"expiry"`
}

func (h *UrlHandler) Create(ctx echo.Context) error {
	var req CreateShortUrlRequest
	if err := ctx.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Validation Error")
	}

	// TODO: fix validation error
	if err := ctx.Validate(&req); err != nil {
		fmt.Println(">>> validation err: ", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Validation Error")
	}

	key, keyErr := h.KeyDb.GetUnused(ctx.Request().Context())
	if keyErr != nil {
		return keyErr
	}

	urlData := &urls.UrlData{Key: key, Url: req.Url, Expiry: req.Expiry}
	if urlErr := h.UrlDb.Insert(ctx.Request().Context(), urlData); urlErr != nil {
		return urlErr
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"url": h.ApiDomain + key,
	})
}
