package config

import (
	"errors"
	"os"
	"strconv"
)

type Config struct {
	DatabaseURL        string
	ServerPort         int
	GoogleClientID     string
	GoogleClientSecret string
	GoogleCallbackURL  string
	WebFrontendURL     string
	JWTSecret          string
	CookieSecure       bool
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

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, errors.New("JWT_SECRET is not set")
	}

	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	if googleClientID == "" {
		return nil, errors.New("GOOGLE_CLIENT_ID is not set")
	}

	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	if googleClientSecret == "" {
		return nil, errors.New("GOOGLE_CLIENT_SECRET is not set")
	}

	return &Config{
		DatabaseURL:        dbURL,
		ServerPort:         port,
		GoogleClientID:     googleClientID,
		GoogleClientSecret: googleClientSecret,
		GoogleCallbackURL:  getEnvOrDefault("GOOGLE_CALLBACK_URL", "http://localhost:8080/auth/callback"),
		WebFrontendURL:     getEnvOrDefault("WEB_FRONTEND_URL", "http://localhost:3000"),
		JWTSecret:          jwtSecret,
		CookieSecure:       os.Getenv("COOKIE_SECURE") == "true",
	}, nil
}

func getEnvOrDefault(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
