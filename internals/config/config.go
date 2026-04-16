package config

import (
	_ "embed"
	"log"

	"go.yaml.in/yaml/v2"
)

//go:embed config.yaml
var configFile []byte

func Load() Config {
	var cfg Config
	if err := yaml.Unmarshal(configFile, &cfg); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	return cfg
}

type Config struct {
	Server  ServerConfig  `yaml:"server"`
	Metrics MetricsConfig `yaml:"metrics"`
}

type ServerConfig struct {
	Port       int   `yaml:"port"`
	Workers    int   `yaml:"workers"`
	QueueSize  int64 `yaml:"queue_size"`
	TokenRate  int64 `yaml:"token_rate"`
	TokenLimit int64 `yaml:"token_limit"`
}

type MetricsConfig struct {
	Port int    `yaml:"port"`
	Path string `yaml:"path"`
}
