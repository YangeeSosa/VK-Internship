package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	GRPCPort string `yaml:"grpc_port"`
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	if err := yaml.NewDecoder(file).Decode(&config); err != nil {
		return nil, err
	}

	// Значения по умолчанию
	if config.GRPCPort == "" {
		config.GRPCPort = ":50051"
	}

	return &config, nil
}
