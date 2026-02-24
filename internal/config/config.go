package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Token string `yaml:"token"`
}

func Load() (*Config, error) {
	// 1. Check environment variable
	if token := os.Getenv("NOTION_API_TOKEN"); token != "" {
		return &Config{Token: token}, nil
	}

	// 2. Check config file
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("could not find home directory: %w", err)
	}

	configPath := filepath.Join(home, ".notion-tui", "config.yaml")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("NOTION_API_TOKEN not set and no config at %s", configPath)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("invalid config file: %w", err)
	}

	if cfg.Token == "" {
		return nil, fmt.Errorf("token is empty in %s", configPath)
	}

	return &cfg, nil
}
