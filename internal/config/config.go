package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

const (
	configPath = "./configs/config.yaml"
)

type PGConfig interface {
	GetDSN() string
}

type HTTPConfig interface {
	GetPort() string
	GetHost() string
	GetTimeout() time.Duration
	GetIdleTimeout() time.Duration
}

type JWTConfig interface {
	GetSecret() string
}

func LoadConfig() (string, error) {
	if err := LoadEnv(); err != nil {
		return "", err
	}

	cfgPath := configPath

	if _, err := os.Stat(cfgPath); err != nil {
		return "", fmt.Errorf("%s file not found", cfgPath)
	}

	return cfgPath, nil
}

func LoadEnv() error {
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("error loading .env file: %w", err)
	}
	return nil
}
