package urlhandlers

import (
	"net/http"

	"github.com/Marie20767/url-shortener/internal/store/keys"
	"github.com/Marie20767/url-shortener/internal/store/urls"
	"github.com/labstack/echo/v4"
)

type URLHandler struct {
	URLsDb *urls.UrlStore
	KeysDb *keys.KeyStore
}

func (h *URLHandler) CreateURL(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"message": "URL created!",
	})
}
