package trojan

import (
	"fmt"

	"github.com/quadraphony/ghostfy/internal/config"
	"github.com/quadraphony/ghostfy/internal/protocols/registry"
)

type Adapter struct{}

func (Adapter) Name() string {
	return "trojan"
}

func (Adapter) Metadata() registry.ProtocolMetadata {
	return registry.ProtocolMetadata{
		Name:         "trojan",
		Category:     "modern-proxy-core",
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
	if cfg.Password == "" {
		return fmt.Errorf("password is required")
	}
	if cfg.TLS == nil || !cfg.TLS.Enabled {
		return fmt.Errorf("tls.enabled must be true for trojan")
	}
	return nil
}

func (Adapter) Build(cfg config.OutboundConfig) (map[string]any, error) {
	out := map[string]any{
		"type":        "trojan",
		"tag":         "proxy",
		"server":      cfg.Server,
		"server_port": cfg.Port,
		"password":    cfg.Password,
	}

	if cfg.TLS != nil && cfg.TLS.Enabled {
		tls := map[string]any{"enabled": true}
		if cfg.TLS.ServerName != "" {
			tls["server_name"] = cfg.TLS.ServerName
		}
		if cfg.TLS.Reality != nil && cfg.TLS.Reality.Enabled {
			tls["reality"] = map[string]any{
				"enabled":    true,
				"public_key": cfg.TLS.Reality.PublicKey,
				"short_id":   cfg.TLS.Reality.ShortID,
			}
			if cfg.TLS.Reality.Fingerprint != "" {
				tls["utls"] = map[string]any{
					"enabled":     true,
					"fingerprint": cfg.TLS.Reality.Fingerprint,
				}
			}
		}
		out["tls"] = tls
	}

	return out, nil
}
