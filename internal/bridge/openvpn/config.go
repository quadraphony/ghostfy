package openvpn

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Profile       string   `json:"profile"`
	LogLevel      string   `json:"log_level"`
	OpenVPNConfig string   `json:"openvpn_config"`
	Executor      string   `json:"openvpn_path,omitempty"`
	Args          []string `json:"args,omitempty"`
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
		c.Profile = "openvpn-bridge"
	}
	if c.LogLevel == "" {
		c.LogLevel = "info"
	}
	c.OpenVPNConfig = strings.TrimSpace(c.OpenVPNConfig)
}

func (c Config) validate() error {
	if c.OpenVPNConfig == "" {
		return errors.New("openvpn_config is required")
	}
	return nil
}
