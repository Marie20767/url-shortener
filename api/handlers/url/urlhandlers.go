package urlhandlers

import (
	"net/http"
	"os"

	"github.com/Marie20767/url-shortener/internal/store/keys"
	"github.com/Marie20767/url-shortener/internal/store/urls"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

type UrlHandler struct {
	UrlDb *urls.UrlStore
	KeyDb *keys.KeyStore
}

func (h *UrlHandler) CreateShortUrl(ctx echo.Context) error {
	if err := godotenv.Load(); err != nil {
		return err
	}

	domain := os.Getenv("API_DOMAIN")
	key := "123xbcaa"
	shortUrl := domain + key

	return ctx.JSON(http.StatusOK, map[string]string{
		"url": shortUrl,
	})
}
