package routes

import (
	"github.com/Marie20767/url-shortener/api/handlers/url"
	"github.com/Marie20767/url-shortener/internal/store/keys"
	"github.com/Marie20767/url-shortener/internal/store/urls"
	"github.com/labstack/echo/v4"
)

func RegisterAll(e *echo.Echo, KeysDb *keys.KeyStore, URLsDb *urls.UrlStore) {
	URLHandler := &urlhandlers.URLHandler{KeysDb: KeysDb, URLsDb: URLsDb}

	e.POST("/create", URLHandler.CreateShortURL)
}
