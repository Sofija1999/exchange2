package config

import (
	"path/filepath"

	"github.com/Bloxico/exchange-gateway/sofija/log"
)

type Environment string

const (
	EnvLocal Environment = "local"
	EnvProd  Environment = "prod"
)

type Config struct {
	Env Environment `yaml:"env"`

	Http     ServerConfig   `yaml:"http" mapstructure:"http"`
	Database DatabaseConfig `yaml:"db" mapstructure:"db"`

	SentryDSN  string `yaml:"sentry_dsn"`
	BaseDomain string `yaml:"base_domain"`

	// ConfigDir from where the config was loaded
	ConfigDir string
}

type ServerConfig struct {
	Port int `yaml:"port"`

	Logger log.Logger
}

type DatabaseConfig struct {
	Name   string `yaml:"name"`
	Host   string `yaml:"host"`
	Port   int    `yaml:"port"`
	User   string `yaml:"user"`
	Pass   string `yaml:"pass"`
	Schema string `yaml:"schema"`
}

func NewConfig(p Provider) (Config, error) {
	var cfg Config

	configFile, err := p.Load(&cfg)
	if err != nil {
		return Config{}, err
	}
	cfg.ConfigDir = filepath.Dir(configFile)

	return cfg, nil
}
