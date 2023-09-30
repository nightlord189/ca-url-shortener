package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	LogLevel string      `yaml:"log_level" env:"LOG_LEVEL"` // trace, debug, info, warn, error, fatal, panic, disabled
	Host     string      `yaml:"host" env:"HOST"`
	HTTP     HTTPConfig  `yaml:"http"`
	Mongo    MongoConfig `yaml:"mongo"`
}

type HTTPConfig struct {
	Port             int    `yaml:"port" env:"HTTP_PORT"`
	AuthTokenExpTime int    `yaml:"auth_token_exp_time" env:"AUTH_TOKEN_EXP_TIME"`
	AuthSecret       string `yaml:"auth_secret" env:"AUTH_SECRET"`
}

type MongoConfig struct {
	Host     string `yaml:"host" env:"MONGO_HOST"`
	Port     int    `yaml:"port" env:"MONGO_PORT"`
	User     string `yaml:"user" env:"MONGO_USER"`
	Password string `yaml:"password" env:"MONGO_PASSWORD"`
	Database string `yaml:"database" env:"MONGO_DATABASE"`
}

func LoadConfig(configFilePath string) (Config, error) {
	var cfg Config
	err := cleanenv.ReadConfig(configFilePath, &cfg)
	if err != nil {
		return cfg, fmt.Errorf("error load config: %w", err)
	}
	return cfg, nil
}
