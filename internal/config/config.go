package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	Port       string `env:"PORT" env-default:"8081"`
	DBHost     string `env:"DB_HOST" env-default:"localhost"`
	DBPort     string `env:"DB_PORT" env-default:"5432"`
	DBUser     string `env:"DB_USER" env-default:"postgres"`
	DBPassword string `env:"PORT" env-default:"xhPgIhZ4"`
	DBName     string `env:"DB_NAME" env-default:"songs_db"`
	APIBaseURL string `env:"API_BASE_URL" env-default:"http://localhost:8081"`
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	return &Config{
		Port:       getEnv("PORT", "8081"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "xhPgIhZ4"),
		DBName:     getEnv("DB_NAME", "songs_db"),
		APIBaseURL: getEnv("API_BASE_URL", "http://localhost:8081"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
