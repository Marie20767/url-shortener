package urlhandlers

import (
	"net/http"

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
	urls.UrlData `bson:",inline"`
	Key          struct{} `bson:"-" json:"-" validate:"-"`
}

func (h *UrlHandler) CreateShortUrl(ctx echo.Context) error {
	var req CreateShortUrlRequest
	if err := ctx.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Validation Error")
	}

	if err := ctx.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Validation Error")
	}

	key, keyErr := h.KeyDb.GetUnusedKey(ctx.Request().Context())
	if keyErr != nil {
		return keyErr
	}

	urlData := &urls.UrlData{Key: key, Url: req.Url, Expiry: req.Expiry}
	if urlErr := h.UrlDb.InsertUrlData(ctx.Request().Context(), urlData); urlErr != nil {
		return urlErr
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"url": h.ApiDomain + key,
	})
}
