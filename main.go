package main

import (
	"database/sql"
	"log"
	"os"

	"web-app/api/routes"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, relying on environment variables")
	}
	
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable not set")
	}

	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
			log.Fatal(err)
	}

	defer dbConn.Close()

	err = dbConn.Ping()
	if err != nil {
		log.Fatalf("Failed to ping DB: %v", err)
	}

	log.Println("Connected to the database successfully!")
	
	e := echo.New()
	
	routes.RegisterAll(e, dbConn)
	
	e.Start(":8080")
}
