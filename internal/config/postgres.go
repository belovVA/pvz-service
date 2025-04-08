package config

import (
	"fmt"
	"strconv"

	"github.com/ilyakaznacheev/cleanenv"
)

type pgConfig struct {
	Name     string `yaml:"database_name" env:"DATABASE_NAME" env-required:"true"`
	Host     string `yaml:"database_host" env:"DATABASE_HOST" env-required:"true"`
	Port     string `yaml:"database_port" env:"DATABASE_PORT" env-required:"true"`
	User     string `yaml:"database_user" env:"DATABASE_USER" env-required:"true"`
	SSLMode  string `yaml:"database_ssl_mode" env:"DATABASE_SSL_MODE" env-required:"true"`
	Password string `env:"DATABASE_PASSWORD" env-required:"true"`
}

func PGConfigLoad() (*pgConfig, error) {
	path, err := LoadConfig()
	if err != nil {
		return nil, err
	}

	var pgCfg pgConfig

	// Читаем конфиг-файл и заполняем нашу структуру
	if err := cleanenv.ReadConfig(path, &pgCfg); err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	if _, err := strconv.Atoi(pgCfg.Port); err != nil {
		return nil, fmt.Errorf("invalid database port: %s", err)
	}

	return &pgCfg, nil
}

func (cfg *pgConfig) GetDSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Host,
		cfg.User,
		cfg.Password,
		cfg.Name,
		cfg.Port,
		cfg.SSLMode,
	)
}
