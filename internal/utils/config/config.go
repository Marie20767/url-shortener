package config

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

var logLevelMap = map[string]slog.Level{
	"debug": slog.LevelDebug,
	"info":  slog.LevelInfo,
	"warn":  slog.LevelWarn,
	"error": slog.LevelError,
}

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
	Port     string
	Domain   string
	Key      *Key
	LogLevel slog.Level
	Url      *Url
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
		"LOG_LEVEL":         nil,
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

	logLevel, ok := logLevelMap[*envVars["LOG_LEVEL"]]
	if !ok {
		return nil, errors.New("log level not set")
	}

	cfg := &cfg{
		Domain:   *envVars["API_DOMAIN"],
		Key:      Key,
		LogLevel: logLevel,
		Port:     *envVars["PORT"],
		Url:      Url,
	}

	return cfg, nil
}
