// Package config handles environment variables.
package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// Config contains environment variables.
type Config struct {
	Port                       string `envconfig:"PORT" default:"8000"`
	DatabasePassword           string `envconfig:"DATABASE_PASSWORD" required:"true"`
	DatabaseUser               string `envconfig:"DATABASE_USER" required:"true"`
	DatabaseURL                string `envconfig:"DATABASE_URL" default:"127.0.0.1"`
	DatabasePort               string `envconfig:"DATABASE_PORT" default:"5432"`
	DatabaseDB                 string `envconfig:"DATABASE_DB" default:"postgres"`
	DatabaseOptions            string `envconfig:"DATABASE_OPTIONS" default:"?sslmode=disable"`
	DatabaseMaxConnections     int    `envconfig:"DATABASE_MAX_CONNECTIONS" default:"12"`
	DatabaseMaxIdleConnections int    `envconfig:"DATABASE_MAX_IDLE_CONNECTIONS" default:"3"`
}

// LoadConfig reads environment variables and populates Config.
func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	var c Config

	err := envconfig.Process("", &c)

	return &c, err
}
