package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type cfg struct {
	Port      string
	Domain    string
	KeyDbUrl  string
	UrlDbUrl  string
	UrlDbName string
}

func ParseEnv() (*cfg, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	envVars := map[string]*string{
		"KEY_DB_URL":  nil,
		"URL_DB_URL":  nil,
		"URL_DB_NAME": nil,
		"PORT":        nil,
		"API_DOMAIN":  nil,
	}

	for key := range envVars {
		value := os.Getenv(key)
		if value == "" {
			return nil, errors.New("not all environment variables are set")
		}
		envVars[key] = &value
	}

	cfg := &cfg{
		KeyDbUrl:  *envVars["KEY_DB_URL"],
		UrlDbUrl:  *envVars["URL_DB_URL"],
		UrlDbName: *envVars["URL_DB_NAME"],
		Port:      *envVars["PORT"],
		Domain:    *envVars["API_DOMAIN"],
	}

	return cfg, nil
}
