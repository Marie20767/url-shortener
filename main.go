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

func run() error {
	ctx := context.Background()

	cfg, err := utils.ParseEnv()
	if err != nil {
		return err
	}

	keysDb, err := keys.NewStore(cfg.KeysDbURL)
	if err != nil {
		return err
	}
	defer keysDb.Close()
	log.Println("connected to keys Db successfully!")

	URLsDb, err := urls.NewStore(cfg.URLsDbURL)
	if err != nil {
		return err
	}
	defer URLsDb.Close(ctx)
	log.Println("connected to URLs Db successfully!")

	e := echo.New()
	routes.RegisterAll(e, keysDb, URLsDb)
	return e.Start(":" + cfg.Port)
}

func main() {
	if err := run(); err != nil {
		log.Println("server closed: ", err)
		os.Exit(1)
	}
}
