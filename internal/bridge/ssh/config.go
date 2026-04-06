package ssh

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Profile    string   `json:"profile"`
	LogLevel   string   `json:"log_level"`
	Binary     string   `json:"binary,omitempty"`
	Args       []string `json:"args,omitempty"`
	Connection string   `json:"connection"`
}

func Load(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("decode config: %w", err)
	}

	cfg.normalize()
	if err := cfg.validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (c *Config) normalize() {
	if c.Profile == "" {
		c.Profile = "ssh-bridge"
	}
	if c.LogLevel == "" {
		c.LogLevel = "info"
	}
	c.Connection = strings.TrimSpace(c.Connection)
}

func (c Config) validate() error {
	if c.Connection == "" {
		return errors.New("connection target is required")
	}
	return nil
}
