package config

import (
	"log"
	"os"
	"github.com/joho/godotenv"
)

type Config struct {
	AppPort    string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
}

func Load() *Config {

	// load .env only in dev
	if os.Getenv("ENV") != "production" {
		_ = godotenv.Load()
	}

	cfg := &Config{
		AppPort:    getEnv("APP_PORT", "8083"),
		DBHost:     getEnv("DB_HOST", "postgres"), // 🔥 Docker fix
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBName:     getEnv("DB_NAME", "lms_db"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	log.Println("lms-service config loaded")
	return cfg
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}