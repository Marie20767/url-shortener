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
	CacheUrl     string
	CronSchedule string
	DbUrl        string
	DbName       string
}

type Key struct {
	CacheUrl     string
	CronSchedule string
	DbUrl        string
}

type cfg struct {
	Domain   string
	Key      *Key
	LogLevel slog.Level
	Port     string
	Url      *Url
}

func ParseEnv() (*cfg, error) {
	// ignore error because in production there will be no .env file
	// env vars will be passed at runtime via docker run command/docker compose
	_ = godotenv.Load()

	envVars := map[string]string{
		"API_DOMAIN":        "",
		"KEY_CACHE_URL":     "",
		"KEY_CRON_SCHEDULE": "",
		"KEY_DB_URL":        "",
		"LOG_LEVEL":         "",
		"SERVER_PORT":              "",
		"URL_CACHE_URL":     "",
		"URL_CRON_SCHEDULE": "",
		"URL_DB_NAME":       "",
		"URL_DB_URL":        "",
	}

	for key := range envVars {
		value := os.Getenv(key)
		if value == "" {
			return nil, fmt.Errorf("%s environment variable not set", key)
		}
		envVars[key] = value
	}

	Key := &Key{
		CacheUrl:     envVars["KEY_CACHE_URL"],
		CronSchedule: envVars["KEY_CRON_SCHEDULE"],
		DbUrl:        envVars["KEY_DB_URL"],
	}

	Url := &Url{
		CacheUrl:     envVars["URL_CACHE_URL"],
		CronSchedule: envVars["URL_CRON_SCHEDULE"],
		DbName:       envVars["URL_DB_NAME"],
		DbUrl:        envVars["URL_DB_URL"],
	}

	logLevel, ok := logLevelMap[envVars["LOG_LEVEL"]]
	if !ok {
		return nil, errors.New("LOG_LEVEL should be one of debug|info|warning|error")
	}

	cfg := &cfg{
		Domain:   envVars["API_DOMAIN"],
		Key:      Key,
		LogLevel: logLevel,
		Port:     envVars["SERVER_PORT"],
		Url:      Url,
	}

	return cfg, nil
}
