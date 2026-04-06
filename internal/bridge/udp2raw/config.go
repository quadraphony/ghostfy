package udp2raw

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Profile    string   `json:"profile"`
	LogLevel   string   `json:"log_level"`
	Executable string   `json:"executable,omitempty"`
	Args       []string `json:"args"`
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

	if cfg.Profile == "" {
		cfg.Profile = "udp2raw-bridge"
	}
	if cfg.LogLevel == "" {
		cfg.LogLevel = "info"
	}
	for i, arg := range cfg.Args {
		cfg.Args[i] = strings.TrimSpace(arg)
	}

	if len(cfg.Args) == 0 {
		return Config{}, fmt.Errorf("args are required")
	}

	return cfg, nil
}
