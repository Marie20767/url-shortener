package config

import (
	"errors"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Url struct {
	DbUrl         string
	DbName        string
	DbTimeout     int
	CacheCapacity int
}

type Key struct {
	DbUrl string
}

type cfg struct {
	Port   string
	Domain string
	Key    *Key
	Url    *Url
}

func ParseEnv() (*cfg, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	envVars := map[string]*string{
		"KEY_DB_URL":         nil,
		"URL_DB_URL":         nil,
		"URL_DB_NAME":        nil,
		"URL_DB_TIMEOUT":     nil,
		"URL_CACHE_CAPACITY": nil,
		"PORT":               nil,
		"API_DOMAIN":         nil,
	}

	for key := range envVars {
		value := os.Getenv(key)
		if value == "" {
			return nil, errors.New("not all environment variables are set")
		}
		envVars[key] = &value
	}

	Key := &Key{
		DbUrl: *envVars["KEY_DB_URL"],
	}

	urlDbTimeout, err := strconv.Atoi(*envVars["URL_DB_TIMEOUT"])
	if err != nil {
		return nil, err
	}
	urlCacheCapacity, err := strconv.Atoi(*envVars["URL_CACHE_CAPACITY"])
	if err != nil {
		return nil, err
	}
	Url := &Url{
		DbUrl:         *envVars["URL_DB_URL"],
		DbName:        *envVars["URL_DB_NAME"],
		DbTimeout:     urlDbTimeout,
		CacheCapacity: urlCacheCapacity,
	}

	cfg := &cfg{
		Key:    Key,
		Url:    Url,
		Port:   *envVars["PORT"],
		Domain: *envVars["API_DOMAIN"],
	}

	return cfg, nil
}
