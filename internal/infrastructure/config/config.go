package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	App      AppConfig      `yaml:"app"`
	Log      LogConfig      `yaml:"log"`
	HTTP     HTTPConfig     `yaml:"http"`
	Postgres PostgresConfig `yaml:"postgres"`
	Migrate  MigrateConfig  `yaml:"migrate"`
}

type AppConfig struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

type LogConfig struct {
	Level string `yaml:"level"`
}

type HTTPConfig struct {
	Port           int  `yaml:"port"`
	UsePreforkMode bool `yaml:"use_prefork_mode"`
}

type PostgresConfig struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	DBName       string `yaml:"dbname"`
	SSLMode      string `yaml:"sslmode"`
	TimeZone     string `yaml:"time_zone"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
	MaxOpenConns int    `yaml:"max_open_conns"`
}

type MigrateConfig struct {
	Dir string `yaml:"dir"`
}

func ConfigFile() string {
	env := os.Getenv("APP_ENV")
	if env == "prod" || env == "production" {
		return "configs/config.prod.yaml"
	}
	return "configs/config.yaml"
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file %s: %w", path, err)
	}
	cfg := new(Config)
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config file %s: %w", path, err)
	}
	return cfg, nil
}
