package main

import (
	"fmt"
	"os"
	"time"
)

// =============================================================================
// Configuration Types
// =============================================================================

// Config holds all configuration for the application
type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
}

// DatabaseConfig holds database connection configuration
type DatabaseConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string
	SSLMode  string
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// =============================================================================
// Configuration Loading
// =============================================================================

// LoadConfig loads and validates all configuration from environment variables
func LoadConfig() (*Config, error) {
	cfg := &Config{
		Database: DatabaseConfig{
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Host:     getEnvOrDefault("DB_HOST", "localhost"),
			Port:     getEnvOrDefault("DB_PORT", "5432"),
			Name:     os.Getenv("DB_NAME"),
			SSLMode:  getEnvOrDefault("DB_SSLMODE", "disable"),
		},
		Server: ServerConfig{
			Port:         getEnvOrDefault("SERVER_PORT", "4000"),
			ReadTimeout:  parseDurationOrDefault("SERVER_READ_TIMEOUT", 5*time.Second),
			WriteTimeout: parseDurationOrDefault("SERVER_WRITE_TIMEOUT", 10*time.Second),
			IdleTimeout:  parseDurationOrDefault("SERVER_IDLE_TIMEOUT", time.Minute),
		},
	}

	// Validate required fields
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate checks that all required configuration is present
func (c *Config) Validate() error {
	missing := []string{}

	if c.Database.User == "" {
		missing = append(missing, "DB_USER")
	}
	if c.Database.Password == "" {
		missing = append(missing, "DB_PASSWORD")
	}
	if c.Database.Name == "" {
		missing = append(missing, "DB_NAME")
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required environment variables: %v", missing)
	}

	return nil
}

// =============================================================================
// Configuration Methods
// =============================================================================

// DSN returns the PostgreSQL connection string
func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Name,
		c.SSLMode,
	)
}

// =============================================================================
// Helper Functions
// =============================================================================

// getEnvOrDefault retrieves an environment variable or returns a default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// parseDurationOrDefault parses a duration from env var or returns a default
func parseDurationOrDefault(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
