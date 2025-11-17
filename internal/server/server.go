package server

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	"github.com/Marie20767/url-shortener/api/handlers"
	"github.com/Marie20767/url-shortener/api/routes"
	"github.com/Marie20767/url-shortener/internal/store/keys"
	"github.com/Marie20767/url-shortener/internal/store/urls"
)

const serverTimeout = 10

type customValidator struct {
	validator *validator.Validate
}

func (cv *customValidator) Validate(i any) error {
	return cv.validator.Struct(i)
}

type Server struct {
	echo *echo.Echo
}

func New(keyStore *keys.KeyStore, urlStore *urls.UrlStore, apiDomain string) *Server {
	server := echo.New()
	server.Validator = &customValidator{validator: validator.New()}
	handler := &handlers.Handler{
		KeyStore:  keyStore,
		UrlStore:  urlStore,
		ApiDomain: apiDomain,
	}
	routes.RegisterAll(server, handler)

	return &Server{
		echo: server,
	}
}

func (s *Server) Start(port string) error {
	err := s.echo.Start(":" + port)
	if err != nil && err != http.ErrServerClosed {
		slog.Error("server error", slog.Any("error", err))
		return err
	}

	return nil
}

func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), serverTimeout*time.Second)
	defer cancel()

	return s.echo.Shutdown(ctx)
}
