package socks5

import (
	"fmt"

	"github.com/quadraphony/ghostfy/internal/config"
	"github.com/quadraphony/ghostfy/internal/protocols/registry"
)

type Adapter struct{}

func (Adapter) Name() string {
	return "socks5"
}

func (Adapter) Metadata() registry.ProtocolMetadata {
	return registry.ProtocolMetadata{
		Name:         "socks5",
		Category:     "proxy-basics",
		Status:       "IMPLEMENTED",
		SupportClass: "Class A",
	}
}

func (Adapter) Validate(cfg config.OutboundConfig) error {
	if cfg.Server == "" {
		return fmt.Errorf("server is required")
	}
	if cfg.Port < 1 || cfg.Port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535")
	}
	if cfg.Version != "" && cfg.Version != "4" && cfg.Version != "4a" && cfg.Version != "5" {
		return fmt.Errorf("unsupported SOCKS version %q", cfg.Version)
	}
	return nil
}

func (Adapter) Build(cfg config.OutboundConfig) (map[string]any, error) {
	version := cfg.Version
	if version == "" {
		version = "5"
	}

	out := map[string]any{
		"type":        "socks",
		"tag":         "proxy",
		"server":      cfg.Server,
		"server_port": cfg.Port,
		"version":     version,
	}
	if cfg.Username != "" {
		out["username"] = cfg.Username
	}
	if cfg.Password != "" {
		out["password"] = cfg.Password
	}

	return out, nil
}
