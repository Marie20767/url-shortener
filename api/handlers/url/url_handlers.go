package urlhandlers

import (
	"net/http"

	"github.com/Marie20767/url-shortener/internal/store"
	"github.com/labstack/echo/v4"
)

type URLHandler struct {
	DB *store.Store
}

func (h *URLHandler) CreateURL(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"message": "URL created!",
	})
}
