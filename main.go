package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"web-app/api/routes"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

func connectDB() (*sql.DB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found")
	}

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
	dbConn, err := connectDB()
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer dbConn.Close()

	log.Println("Connected to DB successfully!")

	e := echo.New()
	routes.RegisterAll(e, dbConn)

	e.Start(":8080")
}
