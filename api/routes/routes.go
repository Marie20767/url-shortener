package routes

import (
	"github.com/Marie20767/url-shortener/api/handlers/url"
	"github.com/Marie20767/url-shortener/internal/store"
	"github.com/labstack/echo/v4"
)

func RegisterAll(e *echo.Echo, db *store.Store) {
	URLHandler := &urlhandlers.URLHandler{DB: db}

	e.POST("/url", URLHandler.CreateURL)
}
