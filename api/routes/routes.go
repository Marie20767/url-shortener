package routes

import (
	handlers "github.com/Marie20767/go-web-app-template/api/handlers/user_handler"
	"github.com/Marie20767/go-web-app-template/internal/store"
	"github.com/labstack/echo/v4"
)

func RegisterAll(e *echo.Echo, db *store.Store) {
	userHandler := &handlers.UserHandler{DB: db}

	e.GET("/hello/:name", userHandler.Hello)
	// other routes here...
}