package main

import (
	"log"
	"os"

	"github.com/Marie20767/go-web-app-template/api/routes"
	"github.com/Marie20767/go-web-app-template/internal/store"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found")
	}

	dbURL := os.Getenv("DATABASE_URL")
	port := os.Getenv("PORT")
	
	if dbURL == "" || port == "" {
		log.Fatalf("Not all environment variables are set: %v", err)
	}

	
	db, err := store.NewStore(dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	log.Println("Connected to DB successfully!")

	e := echo.New()
	routes.RegisterAll(e, db)

	if err := e.Start(":" + port); err != nil {
		log.Fatal(err)
	}
}
