package config

import (
	"errors"
	"os"
	"strconv"
)

type Config struct {
	DatabaseURL string
	ServerPort  int
}

func Load() (*Config, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, errors.New("DATABASE_URL is not set")
	}

	port := 8080
	if portStr := os.Getenv("SERVER_PORT"); portStr != "" {
		parsed, err := strconv.Atoi(portStr)
		if err != nil {
			return nil, errors.New("SERVER_PORT must be a valid integer")
		}
		port = parsed
	}

	return &Config{
		DatabaseURL: dbURL,
		ServerPort:  port,
	}, nil
}
