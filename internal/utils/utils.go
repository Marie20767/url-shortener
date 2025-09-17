package utils

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type config struct {
	Port      string
	KeysDbURL string
	URLsDbURL string
}

func ParseEnv() (*config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	keysDbURL := os.Getenv("KEYS_DB_URL")
	URLsDbURL := os.Getenv("URLS_DB_URL")
	port := os.Getenv("PORT")

	if keysDbURL == "" || URLsDbURL == "" || port == "" {
		return nil, errors.New("not all environment variables are set")
	}

	return &config{
		Port:      port,
		KeysDbURL: keysDbURL,
		URLsDbURL: URLsDbURL,
	}, nil
}
