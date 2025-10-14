package main

import (
	"context"
	"log"
	"os"

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

	cfg, err := utils.ParseEnv()
	if err != nil {
		return err
	}

	keyDb, err := keys.NewStore(cfg.KeyDbUrl)
	if err != nil {
		return err
	}
	defer keyDb.Close()
	log.Println("connected to key Db successfully!")

	urlDb, err := urls.NewStore(cfg.UrlDbUrl)
	if err != nil {
		return err
	}
	defer urlDb.Close(ctx)
	log.Println("connected to url Db successfully!")

	e := echo.New()
	routes.RegisterAll(e, keyDb, urlDb)
	return e.Start(":" + cfg.Port)
}
