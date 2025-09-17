package routes

import (
	"github.com/Marie20767/url-shortener/api/handlers/url"
	"github.com/Marie20767/url-shortener/internal/store/keys"
	"github.com/Marie20767/url-shortener/internal/store/urls"
	"github.com/labstack/echo/v4"
)

func RegisterAll(e *echo.Echo, kDb *keys.KeyStore, uDb *urls.UrlStore) {
	URLHandler := &urlhandlers.URLHandler{KeysDb: kDb, URLsDb: uDb}

	e.POST("/create", URLHandler.CreateShortURL)
}
