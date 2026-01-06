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

const (
	// how long the server will wait to read the entire request after the connection is accepted
	readTimeout = 10 * time.Second

	// how long the server has to write the response after reading the request
	writeTimeout = 10 * time.Second

	// how long to keep a keep-alive connection open waiting for the next request
	idleTimeout = 120 * time.Second

	shutdownTimeout = 10 * time.Second
)

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
	e := echo.New()
	e.Validator = &customValidator{validator: validator.New()}
	handler := &handlers.Handler{
		KeyStore:  keyStore,
		UrlStore:  urlStore,
		ApiDomain: apiDomain,
	}

	routes.RegisterAll(e, handler)

	e.Server.ReadTimeout = readTimeout
	e.Server.WriteTimeout = writeTimeout
	e.Server.IdleTimeout = idleTimeout

	return &Server{
		echo: e,
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
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	return s.echo.Shutdown(ctx)
}
