package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Config struct {
	Web struct {
		APIHost         string        `json:"apiHost"`
		DebugHost       string        `json:"debugHost"`
		ReadTimeout     time.Duration `json:"readTimeout"`
		WriteTimeout    time.Duration `json:"writeTimeout"`
		IdleTimeout     time.Duration `json:"idleTimeout"`
		ShutdownTimeout time.Duration `json:"shutdownTimeout"`
	} `json:"web"`
	AI struct {
		APIKey string `json:"apiKey"`
		APIURL string `json:"apiUrl"`
	} `json:"ai"`
}

func LoadConfig(path string) (*Config, error) {
	fmt.Printf("Loading config from: %s\n", path) // Log the path being used
	configFile, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file at %s: %w", path, err)
	}

	wd, _ := os.Getwd()
	fmt.Println("Current working directory:", wd)

	var config Config
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config JSON: %w", err)
	}

	return &config, nil
}
