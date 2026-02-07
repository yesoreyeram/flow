package config

import (
	"os"
)

// Config holds the application configuration
type Config struct {
	Port        string
	Environment string
	CorsOrigins []string
}

// Load loads the configuration from environment variables
func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "development"
	}

	corsOrigins := []string{"http://localhost:3000"}
	if origins := os.Getenv("CORS_ORIGINS"); origins != "" {
		corsOrigins = []string{origins}
	}

	return &Config{
		Port:        port,
		Environment: env,
		CorsOrigins: corsOrigins,
	}
}
