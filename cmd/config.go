package main

import (
	"fmt"
	"github.com/caarlos0/env/v6"
)

type Config struct {
	LogLevel    string `env:"LOG_LEVEL" envDefault:"debug"`
	MetricsPort string `env:"METRICS_PORT" envDefault:"9102"`
	ApiPort     string `env:"API_PORT" envDefault:"5000"`
}

func NewConfig() (*Config, error) {

	cfg := &Config{}

	err := env.Parse(cfg)
	if err != nil {
		return nil, fmt.Errorf("could not parse config. %w", err)
	}

	return cfg, nil
}
