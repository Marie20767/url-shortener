package utils

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type config struct {
	Port     string
	KeyDbUrl string
	UrlDbUrl string
}

func ParseEnv() (*config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	keyDbUrl := os.Getenv("KEY_DB_URL")
	urlDbUrl := os.Getenv("URL_DB_URL")
	port := os.Getenv("PORT")

	if keyDbUrl == "" || urlDbUrl == "" || port == "" {
		return nil, errors.New("not all environment variables are set")
	}

	return &config{
		Port:     port,
		KeyDbUrl: keyDbUrl,
		UrlDbUrl: urlDbUrl,
	}, nil
}
