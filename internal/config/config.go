package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	LogLevel string `yaml:"logLevel" json:"logLevel"`
}

func New(configFile string) (*Config, error) {
	config := Config{}

	f, err := os.Open(configFile)
	if err != nil {
		return nil, fmt.Errorf("config: failed to open config file: %w", err)
	}
	defer f.Close()

	if err := yaml.NewDecoder(f).Decode(&config); err != nil {
		return nil, fmt.Errorf("config: invalid config syntax: %w", err)
	}

	return &config, nil
}
