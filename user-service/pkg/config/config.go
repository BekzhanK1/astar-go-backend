package config

import (
	"log"
	"user-service/internal/utils"

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
		utils.GetEnv("DATABASE_HOST", "localhost"),
		utils.GetEnv("DATABASE_PORT", "5432"),
		utils.GetEnv("DATABASE_USER", "user"),
		utils.GetEnv("DATABASE_PASSWORD", "password"),
		utils.GetEnv("DATABASE_NAME", "dbname"),
		utils.GetEnv("DATABASE_SSLMODE", "disable"),
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
