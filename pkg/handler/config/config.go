package config

import (
	"fmt"
	"os"
	"time"

	"github.com/goccy/go-yaml"
)

type Config struct {
	DB     DB     `yaml:"db"`
	Worker Worker `yaml:"worker"`
	Server Server `yaml:"server"`
}

type DB struct {
	Host       string `yaml:"host"`
	Port       string `yaml:"port"`
	Username   string `yaml:"username"`
	Password   string `yaml:"password"`
	Name       string `yaml:"name"`
	DisableTLS bool   `yaml:"disable_tls"`
}

type Worker struct {
	Name              string        `yaml:"name"`
	RetryCount        int           `yaml:"retry_count"`
	Size              int           `yaml:"size"`
	BackoffMultiplier int           `yaml:"backoff_multiplier"`
	DelayInterval     time.Duration `yaml:"delay_interval"`
}

type Server struct {
	Addr string `yaml:"addr"`
}

func New(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading file %s: %w", path, err)
	}
	cfg := &Config{}
	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling data: %w", err)
	}

	return cfg, nil
}
