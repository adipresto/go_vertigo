package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Database struct {
		Path string `yaml:"path"`
	} `yaml:"database"`
	Network struct {
		Centrifugo struct {
			Enabled bool   `yaml:"enabled"`
			URL     string `yaml:"url"`
		} `yaml:"centrifugo"`
		MQTT struct {
			Enabled bool   `yaml:"enabled"`
			URL     string `yaml:"url"`
		} `yaml:"mqtt"`
	} `yaml:"network"`
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`
}

func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
