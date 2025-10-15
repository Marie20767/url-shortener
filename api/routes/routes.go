package routes

import (
	"github.com/Marie20767/url-shortener/api/handlers/url"
	"github.com/labstack/echo/v4"
)

func RegisterAll(e *echo.Echo, urlHandler *urlhandlers.UrlHandler) {
	e.POST("/create", urlHandler.CreateShortUrl)
}
