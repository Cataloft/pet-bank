package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env         string `yaml:"env"`
	PostgresURL string `yaml:"postgres_url"`
	HTTPServer  `yaml:"http_server"`
}

type HTTPServer struct {
	Address              string        `yaml:"address"`
	AccessTokenDuration  time.Duration `yaml:"access_token_duration"`
	RefreshTokenDuration time.Duration `yaml:"refresh_token_duration"`
	Secret               string        `yaml:"secret"`
}

func MustLoad() *Config {
	configPath := os.Getenv("config_path")
	if configPath == "" {
		log.Fatal("config path is not set")
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatal("cannot read cfg")
	}

	return &cfg
}
