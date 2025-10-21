package urlhandlers

import (
	"net/http"

	"github.com/Marie20767/url-shortener/internal/store/keys"
	"github.com/Marie20767/url-shortener/internal/store/urls"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

type UrlHandler struct {
	UrlDb     *urls.UrlStore
	KeyDb     *keys.KeyStore
	ApiDomain string
}

func (h *UrlHandler) Create(ctx echo.Context) error {
	if err := godotenv.Load(); err != nil {
		return err
	}

	key := "123xbcaa"
	shortUrl := h.ApiDomain + key

	return ctx.JSON(http.StatusOK, map[string]string{
		"url": shortUrl,
	})
}
