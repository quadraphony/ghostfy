package trusttunnel

import (
	"fmt"

	"github.com/quadraphony/ghostfy/internal/config"
	"github.com/quadraphony/ghostfy/internal/protocols/registry"
)

type Adapter struct{}

func (Adapter) Name() string {
	return "trusttunnel"
}

func (Adapter) Metadata() registry.ProtocolMetadata {
	return registry.ProtocolMetadata{
		Name:         "trusttunnel",
		Category:     "special",
		Status:       "PLANNED",
		SupportClass: "Class C",
	}
}

func (Adapter) Validate(cfg config.OutboundConfig) error {
	if cfg.Server == "" {
		return fmt.Errorf("server is required")
	}
	if cfg.Port < 1 || cfg.Port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535")
	}
	return nil
}

func (Adapter) Build(cfg config.OutboundConfig) (map[string]any, error) {
	return map[string]any{
		"type":        "trusttunnel",
		"tag":         "proxy",
		"server":      cfg.Server,
		"server_port": cfg.Port,
		"password":    cfg.Password,
		"mode":        cfg.Flow,
	}, nil
}
