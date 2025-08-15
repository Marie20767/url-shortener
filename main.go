package main

import (
	"errors"
	"log"
	"os"

	"github.com/Marie20767/go-web-app-template/api/routes"
	"github.com/Marie20767/go-web-app-template/internal/store"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

func run() error {
	if err := godotenv.Load(); err != nil {
		return err
	}

	dbURL := os.Getenv("DATABASE_URL")
	port := os.Getenv("PORT")
	if dbURL == "" || port == "" {
		return errors.New("not all environment variables are set")
	}

	db, err := store.NewStore(dbURL)
	if err != nil {
		return err
	}
	defer db.Close()
	log.Println("connected to DB successfully!")

	e := echo.New()
	routes.RegisterAll(e, db)
	return e.Start(":" + port);
}

func main() {
	if err := run(); err != nil {
		log.Println("server closed: ", err)
		os.Exit(1)
	}
}
