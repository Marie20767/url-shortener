package main

import (
	"context"
	"log"
	"os"

	"github.com/Marie20767/url-shortener/api/handlers/url"
	"github.com/Marie20767/url-shortener/api/routes"
	"github.com/Marie20767/url-shortener/internal/keygenerator"
	"github.com/Marie20767/url-shortener/internal/store/keys"
	"github.com/Marie20767/url-shortener/internal/store/urls"
	"github.com/Marie20767/url-shortener/internal/utils/config"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i any) error {
	return cv.validator.Struct(i)
}

func main() {
	if err := run(); err != nil {
		log.Println("server closed: ", err)
		os.Exit(1)
	}
}

func run() error {
	ctx := context.Background()

	cfg, err := config.ParseEnv()
	if err != nil {
		return err
	}

	keyStore, err := keys.New(ctx, cfg.Key)
	if err != nil {
		return err
	}
	defer keyStore.Close()
	log.Println("connected to key db successfully!")

	urlStore, err := urls.New(cfg.Url)
	if err != nil {
		return err
	}
	defer urlStore.Close(ctx)
	log.Println("connected to url db successfully!")

	// TODO: change to only generating keys in url request handler if no more keys available
	keyGen := keygenerator.New(keyStore)
	keyGenErr := keyGen.Generate(ctx)
	if keyGenErr != nil {
		return keyGenErr
	}
	log.Println("generated url keys successfully!")

	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	urlHandler := &urlhandlers.UrlHandler{
		KeyStore:  keyStore,
		UrlStore:  urlStore,
		ApiDomain: cfg.Domain,
	}
	routes.RegisterAll(e, urlHandler)
	return e.Start(":" + cfg.Port)
}
