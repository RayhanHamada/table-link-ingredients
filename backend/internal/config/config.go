package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application.
type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Server   ServerConfig
}

// AppConfig holds application-level configuration.
type AppConfig struct {
	Name string
	Env  string
}

// DatabaseConfig holds database connection configuration.
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	Schema   string
	PoolMax  int
}

// DSN returns the PostgreSQL connection string.
func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable&search_path=%s",
		d.User, d.Password, d.Host, d.Port, d.DBName, d.Schema,
	)
}

// ServerConfig holds HTTP server configuration.
type ServerConfig struct {
	Port int
}

// Load reads configuration from environment variables.
// It attempts to load a .env file first, but does not fail if one is absent.
func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		App: AppConfig{
			Name: getEnv("APP_NAME", "tablelink-backend"),
			Env:  getEnv("APP_ENV", "development"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "tablelink"),
			Schema:   getEnv("DB_SCHEMA", "public"),
			PoolMax:  getEnvInt("DB_POOL_MAX", 20),
		},
		Server: ServerConfig{
			Port: getEnvInt("SERVER_PORT", 3000),
		},
	}

	return cfg, nil
}

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	raw := os.Getenv(key)
	if raw == "" {
		return fallback
	}
	val, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return val
}
