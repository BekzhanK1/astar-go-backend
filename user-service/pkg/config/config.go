package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseDSN string
}

// LoadConfig loads configuration from environment variables or .env file.
func LoadConfig() *Config {
	// Load .env file (optional: can skip in production if env vars are directly set)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, loading environment variables.")
	}

	dsn := buildDSN(
		getEnv("DATABASE_HOST", "localhost"),
		getEnv("DATABASE_PORT", "5432"),
		getEnv("DATABASE_USER", "user"),
		getEnv("DATABASE_PASSWORD", "password"),
		getEnv("DATABASE_NAME", "dbname"),
		getEnv("DATABASE_SSLMODE", "disable"),
	)

	return &Config{DatabaseDSN: dsn}
}

func buildDSN(host, port, user, password, dbname, sslmode string) string {
	return "host=" + host +
		" port=" + port +
		" user=" + user +
		" password=" + password +
		" dbname=" + dbname +
		" sslmode=" + sslmode
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
