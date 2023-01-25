package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Config struct {
	Grpc struct {
		Host string `json:"host"`
		Port string `json:"port"`
		Type string `json:"type"`
	} `json:"grpc"`
	Database struct {
		Host     string `json:"host"`
		Port     string `json:"port"`
		Username string `json:"username"`
		Password string `json:"password"`
		Name     string `json:"name"`
		SSLMode  string `json:"ssl_mode"`
	} `json:"database"`
}

func InitConfig(path string) (*Config, error) {
	log.Println("init config")
	var cfg Config
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("init configs: %w", err)
	}
	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("init configs: %w", err)
	}
	return &cfg, nil
}
