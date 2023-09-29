package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	LogLevel string     `yaml:"log_level" env:"LOG_LEVEL"` // trace, debug, info, warn, error, fatal, panic, disabled
	Host     string     `yaml:"host" env:"HOST"`
	HTTP     HTTPConfig `yaml:"http"`
}

type HTTPConfig struct {
	Port             int    `yaml:"port" env:"HTTP_PORT"`
	AuthTokenExpTime int    `yaml:"auth_token_exp_time" env:"AUTH_TOKEN_EXP_TIME"`
	AuthSecret       string `yaml:"auth_secret" env:"AUTH_SECRET"`
}

func LoadConfig(configFilePath string) (Config, error) {
	var cfg Config
	err := cleanenv.ReadConfig(configFilePath, &cfg)
	if err != nil {
		return cfg, fmt.Errorf("error load config: %w", err)
	}
	return cfg, nil
}
