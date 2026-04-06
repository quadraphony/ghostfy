package vmess

import (
	"fmt"

	"github.com/quadraphony/ghostfy/internal/config"
	"github.com/quadraphony/ghostfy/internal/protocols/registry"
)

type Adapter struct{}

func (Adapter) Name() string {
	return "vmess"
}

func (Adapter) Metadata() registry.ProtocolMetadata {
	return registry.ProtocolMetadata{
		Name:         "vmess",
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
	if cfg.UUID == "" {
		return fmt.Errorf("uuid is required")
	}
	if cfg.Security == "" {
		return fmt.Errorf("security is required")
	}
	return nil
}

func (Adapter) Build(cfg config.OutboundConfig) (map[string]any, error) {
	out := map[string]any{
		"type":        "vmess",
		"tag":         "proxy",
		"server":      cfg.Server,
		"server_port": cfg.Port,
		"uuid":        cfg.UUID,
		"alter_id":    cfg.AlterID,
		"security":    cfg.Security,
	}

	if cfg.Network != "" {
		out["network"] = cfg.Network
	}

	if cfg.Transport != nil && cfg.Transport.Type != "" && cfg.Transport.Type != "tcp" {
		trans := map[string]any{
			"type": cfg.Transport.Type,
		}
		if cfg.Transport.Path != "" {
			trans["path"] = cfg.Transport.Path
		}
		if cfg.Transport.Host != "" {
			trans["host"] = cfg.Transport.Host
		}
		out["transport"] = trans
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
