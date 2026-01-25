package config

import (
	"fmt"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	APP_ENV    string `env:"APP_ENV" envDefault:"prod"`
	ServerPort string `env:"SERVER_PORT" envDefault:":8080"`
	DBHost     string `env:"DB_HOST" envDefault:"localhost"`
	DBPort     string `env:"DB_PORT" envDefault:"5432"`
	DBName     string `env:"DB_NAME" envDefault:"shorturl"`
	DBUser     string `env:"DB_USER" envDefault:"arjun"`
	DBPassword string `env:"DB_PASSWORD" envDefault:""`
}

func LoadConfig() (Config, error) {
	// Determine the environment
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		appEnv = "prod"
	}
	// Load .env.local if dev
	if appEnv == "dev" {
		_ = godotenv.Load(".env.local")
	}
	// Parse into struct ONCE
	cfg, err := env.ParseAs[Config]()
	if err != nil {
		return Config{}, err
	}
	fmt.Println(cfg)
	return cfg, nil
}
