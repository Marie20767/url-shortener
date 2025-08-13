package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Marie20767/go-web-app-template/internal/store"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	DB *store.Store
}

func (h *UserHandler) Hello(c echo.Context) error {
	name := c.Param("name")

	user, err := h.DB.Queries.GetUserByName(context.Background(), name)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "User not found"})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": fmt.Sprintf("Hello, %s!", user.Name),
	})
}
