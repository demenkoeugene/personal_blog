package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Port int
	Auth struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"auth"`
}

func GetConfigPath() string {
	return "config/config.yaml"
}

func (c *Config) GetPort() int {
	return 8080
}

func Load(path string) (*Config, error) {
	config := &Config{}

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error opening config: %w", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(config); err != nil {
		return nil, fmt.Errorf("error parsing config: %w", err)
	}

	return config, nil
}
