package config

import "os"

type Config struct {
	Port           string
	DatabaseURL    string
	OpenSearchURL  string
}

func Load() (*Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://marketplace:marketplace_dev_secret@localhost:5432/marketplace?sslmode=disable"
	}
	osURL := os.Getenv("OPENSEARCH_URL")
	if osURL == "" {
		osURL = "http://localhost:9200"
	}
	return &Config{
		Port:          port,
		DatabaseURL:   dbURL,
		OpenSearchURL: osURL,
	}, nil
}
