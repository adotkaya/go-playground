package main

import (
	"fmt"
	"log"
	"os"
)

func loadDatabaseConfig(errorLog *log.Logger) string {
	config := map[string]string{
		"DB_USER":     os.Getenv("DB_USER"),
		"DB_PASSWORD": os.Getenv("DB_PASSWORD"),
		"DB_HOST":     os.Getenv("DB_HOST"),
		"DB_PORT":     os.Getenv("DB_PORT"),
		"DB_NAME":     os.Getenv("DB_NAME"),
		"DB_SSLMODE":  os.Getenv("DB_SSLMODE"),
	}

	// Check for missing variables
	missing := []string{}
	for key, value := range config {
		if value == "" {
			missing = append(missing, key)
		}
	}

	if len(missing) > 0 {
		errorLog.Fatalf("Missing required environment variables: %v", missing)
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		config["DB_USER"],
		config["DB_PASSWORD"],
		config["DB_HOST"],
		config["DB_PORT"],
		config["DB_NAME"],
		config["DB_SSLMODE"],
	)
}
