package config

import (
	"os"
)

type Config struct {
	GatewayHost string `env:"GATEWAY_HOST"`
	GatewayPort string `env:"GATEWAY_PORT"`
	AuthHost    string `env:"AUTH_HOST"`
	AuthPort    string `env:"AUTH_PORT"`
}

func Load() *Config {
	cfg := &Config{}

	cfg.GatewayHost = os.Getenv("GATEWAY_HOST")
	cfg.GatewayPort = os.Getenv("GATEWAY_PORT")

	cfg.AuthHost = os.Getenv("AUTH_HOST")
	cfg.AuthPort = os.Getenv("AUTH_PORT")

	return cfg
}
