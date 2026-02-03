package config

import "os"

type Config struct {
	Port        string
	DatabaseURL string
	AdminAPIKey string
}

func Load() (*Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8086"
	}
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://marketplace:marketplace_dev_secret@localhost:5432/marketplace?sslmode=disable"
	}
	apiKey := os.Getenv("ADMIN_API_KEY")
	if apiKey == "" {
		apiKey = "dev_admin_secret_change_in_production"
	}
	return &Config{
		Port:        port,
		DatabaseURL: dbURL,
		AdminAPIKey: apiKey,
	}, nil
}
