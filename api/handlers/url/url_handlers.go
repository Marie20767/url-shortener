package urlhandlers

import (
	"net/http"
	"os"

	"github.com/Marie20767/url-shortener/internal/store/keys"
	"github.com/Marie20767/url-shortener/internal/store/urls"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

type URLHandler struct {
	URLsDb *urls.UrlStore
	KeysDb *keys.KeyStore
}


func (h *URLHandler) CreateShortURL(c echo.Context) error {
	if err := godotenv.Load(); err != nil {
		return err
	}

	domain := os.Getenv("API_DOMAIN")
	key := "123xbcaa"
	shortURL := domain + key

	return c.JSON(http.StatusOK, map[string]string{
		"url": shortURL,
	})
}
