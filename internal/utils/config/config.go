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
		"URL_CRON_SCHEDULE": nil,
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
		CacheUrl:     *envVars["KEY_CACHE_URL"],
		CronSchedule: *envVars["KEY_CRON_SCHEDULE"],
		DbUrl:        *envVars["KEY_DB_URL"],
	}
	Url := &Url{
		CacheUrl:     *envVars["URL_CACHE_URL"],
		CronSchedule: *envVars["URL_CRON_SCHEDULE"],
		DbName:       *envVars["URL_DB_NAME"],
		DbUrl:        *envVars["URL_DB_URL"],
	}

	cfg := &cfg{
		Domain:   *envVars["API_DOMAIN"],
		Key:      Key,
		LogLevel: logLevelMap[*envVars["LOG_LEVEL"]],
		Port:     *envVars["PORT"],
		Url:      Url,
	}

	return cfg, nil
}
