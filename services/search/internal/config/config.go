package config

import "os"

type Config struct {
	Port          string
	OpenSearchURL string
}

func Load() (*Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}
	osURL := os.Getenv("OPENSEARCH_URL")
	if osURL == "" {
		osURL = "http://localhost:9200"
	}
	return &Config{
		Port:          port,
		OpenSearchURL: osURL,
	}, nil
}
