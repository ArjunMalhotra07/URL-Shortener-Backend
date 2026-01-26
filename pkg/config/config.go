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
	DBDSN      string `env:"DB_DSN" envDefault:""`
	DBMinConns int32  `env:"DB_MIN_CONNS" envDefault:"1"`
	DBMaxConns int32  `env:"DB_MAX_CONNS" envDefault:"10"`
	// Lifetime caps how long a connection can live; zero means unlimited.
	DBMaxConnLifetimeSeconds int32 `env:"DB_MAX_CONN_LIFETIME_SECONDS" envDefault:"0"`
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
