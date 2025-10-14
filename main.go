package main

import (
	"context"
	"log"
	"os"

	"github.com/Marie20767/url-shortener/api/handlers/url"
	"github.com/Marie20767/url-shortener/api/routes"
	"github.com/Marie20767/url-shortener/internal/store/keys"
	"github.com/Marie20767/url-shortener/internal/store/urls"
	"github.com/Marie20767/url-shortener/internal/utils"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

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

	keyDb, err := keys.NewStore(cfg.KeyDbUrl)
	if err != nil {
		return err
	}
	defer keyDb.Close()
	log.Println("connected to key Db successfully!")

	urlDb, err := urls.NewStore(cfg.UrlDbUrl, cfg.UrlDbName)
	if err != nil {
		return err
	}
	defer urlDb.Close(ctx)
	log.Println("connected to url Db successfully!")

	e := echo.New()
	urlHandler := &urlhandlers.UrlHandler{KeyDb: keyDb, UrlDb: urlDb, ApiDomain: cfg.Domain}
	routes.RegisterAll(e, urlHandler)
	return e.Start(":" + cfg.Port)
}
