package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type jwtConfig struct {
	Jwt string `env:"JWT_SECRET" env-required:"true"`
}

func (j *jwtConfig) GetSecret() string {
	return j.Jwt
}

func JWTConfigLoad() (*jwtConfig, error) {
	var jwtCfg jwtConfig

	if err := cleanenv.ReadEnv(&jwtCfg); err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	if len(jwtCfg.Jwt) == 0 {
		return nil, fmt.Errorf("JWT_SECRET enviroment doesnt exist")
	}

	return &jwtCfg, nil

}
