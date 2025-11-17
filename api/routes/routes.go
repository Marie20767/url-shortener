package routes

import (
	"github.com/labstack/echo/v4"

	"github.com/Marie20767/url-shortener/api/handlers"
)

func RegisterAll(e *echo.Echo, h *handlers.Handler) {
	e.GET("/health", h.HealthCheck)
	e.POST("/create", h.CreateShort)
	e.GET("/:key", h.GetLong)
}
