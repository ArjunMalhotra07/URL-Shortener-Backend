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
	DBMaxConns int32  `env:"DB_MAX_CONNS" envDefault:"20"`
	// MaxConnLifetime caps how long a connection can live before being closed and replaced.
	DBMaxConnLifetimeSeconds int32 `env:"DB_MAX_CONN_LIFETIME_SECONDS" envDefault:"1800"`
	// MaxConnIdleTime caps how long an idle connection stays in the pool.
	DBMaxConnIdleTimeSeconds int32 `env:"DB_MAX_CONN_IDLE_TIME_SECONDS" envDefault:"300"`

	// JWT settings
	JWTSecret                 string `env:"JWT_SECRET" envDefault:""`
	JWTAccessTokenExpiryMins  int    `env:"JWT_ACCESS_TOKEN_EXPIRY_MINS" envDefault:"60"`
	JWTRefreshTokenExpiryDays int    `env:"JWT_REFRESH_TOKEN_EXPIRY_DAYS" envDefault:"30"`

	// Frontend URL for redirects
	FrontendURL string `env:"FRONTEND_URL" envDefault:"http://localhost:5173"`
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
