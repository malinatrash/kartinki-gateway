package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	GatewayHost string `env:"GATEWAY_HOST"`
	GatewayPort string `env:"GATEWAY_PORT"`
	AuthHost    string `env:"AUTH_HOST"`
	AuthPort    string `env:"AUTH_PORT"`
	JWTSecret   string `env:"JWT_SECRET"`
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	cfg := &Config{}

	cfg.GatewayHost = os.Getenv("GATEWAY_HOST")
	cfg.GatewayPort = os.Getenv("GATEWAY_PORT")
	cfg.AuthHost = os.Getenv("AUTH_HOST")
	cfg.AuthPort = os.Getenv("AUTH_PORT")
	cfg.JWTSecret = os.Getenv("JWT_SECRET")

	return cfg
}
