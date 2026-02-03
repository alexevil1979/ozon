// Package config — env-based config. Feature flags via env later.
package config

import "os"

// Config holds service config. All from env for 12-factor.
type Config struct {
	Port          string
	DatabaseURL   string
	JWTSecret     string
	JWTExpiry     string // e.g. "24h"
	RefreshExpiry string // e.g. "168h" (7 days)
}

// Load reads config from environment. Production: use explicit env vars or secret manager.
func Load() (*Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://marketplace:marketplace_dev_secret@localhost:5432/marketplace?sslmode=disable"
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "dev_jwt_secret_change_in_production"
	}
	jwtExpiry := os.Getenv("JWT_EXPIRY")
	if jwtExpiry == "" {
		jwtExpiry = "24h"
	}
	refreshExpiry := os.Getenv("REFRESH_EXPIRY")
	if refreshExpiry == "" {
		refreshExpiry = "168h" // 7 days
	}
	return &Config{
		Port:          port,
		DatabaseURL:   dbURL,
		JWTSecret:     jwtSecret,
		JWTExpiry:     jwtExpiry,
		RefreshExpiry: refreshExpiry,
	}, nil
}
