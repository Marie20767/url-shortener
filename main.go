package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/Marie20767/go-web-app-template/api/routes"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

func connectDB() (*sql.DB, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable not set")
	}

	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("db connection error: %w", err)
	}

	err = dbConn.Ping()
	if err != nil {
		dbConn.Close()
		return nil, fmt.Errorf("failed to ping DB: %w", err)
	}

	return dbConn, nil
}

func main() {
	err := godotenv.Load()
	port := ":" + os.Getenv("PORT")

	if err != nil {
		log.Println("Warning: .env file not found")
	}

	dbConn, err := connectDB()
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer dbConn.Close()

	log.Println("Connected to DB successfully!")

	e := echo.New()
	routes.RegisterAll(e, dbConn)

	if err := e.Start(port); err != nil {
		log.Fatal(err)
	}
}
