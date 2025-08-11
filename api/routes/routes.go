package routes

import (
	"database/sql"
	"web-app/api/routes/handlers"

	"github.com/labstack/echo/v4"
)

func RegisterAll(e *echo.Echo, db *sql.DB) {
	userHandler := &handlers.UserHandler{DB: db}

	e.GET("/hello/:name", userHandler.Hello)
	// other routes here...
}