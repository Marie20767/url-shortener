package routes

import (
	"github.com/Marie20767/url-shortener/api/handlers/url"
	"github.com/Marie20767/url-shortener/internal/store/keys"
	"github.com/Marie20767/url-shortener/internal/store/urls"
	"github.com/labstack/echo/v4"
)

func RegisterAll(e *echo.Echo, keyDb *keys.KeyStore, urlDb *urls.UrlStore) {
	urlHandler := &urlhandlers.UrlHandler{KeyDb: keyDb, UrlDb: urlDb}

	e.POST("/create", urlHandler.CreateShortUrl)
}
