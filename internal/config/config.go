package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	LogLevel          string            `yaml:"logLevel" json:"logLevel"`
	ProcessesSettings ProcessesSettings `yaml:"processesSettings" json:"processesSettings"`
}

type ProcessesSettings struct {
	Size                int                 `yaml:"size" json:"size"`
	CustomFilterSetting CustomFilterSetting `yaml:"customFilterSetting" json:"customFilterSetting"`
}

type ProcessSetting struct {
	Size int `yaml:"size" json:"size"`
}

type CustomFilterSetting struct {
	ProcessSetting
	MinValue int `yaml:"minValue" json:"minValue"`
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
