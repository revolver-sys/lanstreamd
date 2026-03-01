package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Proxy struct {
		Listen            string `yaml:"listen"`
		ConnectTimeoutMs  int    `yaml:"connect_timeout_ms"`
		IdleConnTimeoutMs int    `yaml:"idle_conn_timeout_ms"`
	} `yaml:"proxy"`

	API struct {
		Listen string `yaml:"listen"`
	} `yaml:"api"`

	Metrics struct {
		Enabled bool `yaml:"enabled"`
	} `yaml:"metrics"`
}

func Load(path string) (*Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var c Config
	if err := yaml.Unmarshal(b, &c); err != nil {
		return nil, fmt.Errorf("parse yaml: %w", err)
	}

	// defaults
	if c.Proxy.Listen == "" {
		c.Proxy.Listen = "127.0.0.1:3128"
	}
	if c.API.Listen == "" {
		c.API.Listen = "127.0.0.1:8080"
	}
	if c.Proxy.ConnectTimeoutMs == 0 {
		c.Proxy.ConnectTimeoutMs = 8000
	}
	if c.Proxy.IdleConnTimeoutMs == 0 {
		c.Proxy.IdleConnTimeoutMs = 30000
	}

	return &c, nil
}
