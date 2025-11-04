package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Url struct {
	DbUrl    string
	DbName   string
	CacheUrl string
}

type Key struct {
	DbUrl        string
	CacheUrl     string
	CronSchedule string
}

type cfg struct {
	Port   string
	Domain string
	Key    *Key
	Url    *Url
}

func ParseEnv() (*cfg, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("failed to parse env vars: %w", err)
	}

	envVars := map[string]*string{
		"API_DOMAIN":        nil,
		"KEY_CACHE_URL":     nil,
		"KEY_CRON_SCHEDULE": nil,
		"KEY_DB_URL":        nil,
		"PORT":              nil,
		"URL_CACHE_URL":     nil,
		"URL_DB_NAME":       nil,
		"URL_DB_URL":        nil,
	}

	for key := range envVars {
		value := os.Getenv(key)
		if value == "" {
			return nil, errors.New("not all environment variables are set")
		}
		envVars[key] = &value
	}

	Key := &Key{
		DbUrl:        *envVars["KEY_DB_URL"],
		CacheUrl:     *envVars["KEY_CACHE_URL"],
		CronSchedule: *envVars["KEY_CRON_SCHEDULE"],
	}
	Url := &Url{
		DbUrl:    *envVars["URL_DB_URL"],
		DbName:   *envVars["URL_DB_NAME"],
		CacheUrl: *envVars["URL_CACHE_URL"],
	}

	cfg := &cfg{
		Key:    Key,
		Url:    Url,
		Port:   *envVars["PORT"],
		Domain: *envVars["API_DOMAIN"],
	}

	return cfg, nil
}
