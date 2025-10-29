package routes

import (
	"github.com/labstack/echo/v4"

	urlhandlers "github.com/Marie20767/url-shortener/api/handlers/url"
)

func RegisterAll(e *echo.Echo, urlHandler *urlhandlers.UrlHandler) {
	e.POST("/create", urlHandler.CreateShort)
	e.GET("/:key", urlHandler.GetLong)
}
