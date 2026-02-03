package config

import "os"

type Config struct {
	Port        string
	CatalogURL  string
	JWTSecret   string
}

func Load() (*Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8084"
	}
	catalogURL := os.Getenv("CATALOG_URL")
	if catalogURL == "" {
		catalogURL = "http://localhost:8082"
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "dev_jwt_secret_change_in_production"
	}
	return &Config{
		Port:       port,
		CatalogURL: catalogURL,
		JWTSecret:  jwtSecret,
	}, nil
}
