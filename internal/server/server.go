package server

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	handlers "github.com/Marie20767/url-shortener/api/handlers/url"
	"github.com/Marie20767/url-shortener/api/routes"
	"github.com/Marie20767/url-shortener/internal/store/keys"
	"github.com/Marie20767/url-shortener/internal/store/urls"
)

type customValidator struct {
	validator *validator.Validate
}

func (cv *customValidator) Validate(i any) error {
	return cv.validator.Struct(i)
}

type Server struct {
	echo *echo.Echo
	port string
}

func New(keyStore *keys.KeyStore, urlStore *urls.UrlStore, apiDomain, port string) *Server {
	server := echo.New()
	server.Validator = &customValidator{validator: validator.New()}
	urlHandler := &handlers.UrlHandler{
		KeyStore:  keyStore,
		UrlStore:  urlStore,
		ApiDomain: apiDomain,
	}
	routes.RegisterAll(server, urlHandler)

	return &Server{
		echo: server,
	}
}

func (s *Server) Start() error {
	err := s.echo.Start(":" + s.port)
	if err != nil && err != http.ErrServerClosed {
		slog.Error("server error", slog.Any("error", err))
		return err
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	err := s.echo.Shutdown(ctx)

	return err
}
