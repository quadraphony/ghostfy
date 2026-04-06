package vless

import (
	"fmt"

	"github.com/quadraphony/ghostfy/internal/config"
	"github.com/quadraphony/ghostfy/internal/protocols/registry"
)

type Adapter struct{}

func (Adapter) Name() string {
	return "vless"
}

func (Adapter) Metadata() registry.ProtocolMetadata {
	return registry.ProtocolMetadata{
		Name:         "vless",
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
	if cfg.TLS == nil || !cfg.TLS.Enabled {
		return fmt.Errorf("tls.enabled must be true")
	}
	if cfg.Flow != "" && cfg.Flow != "xtls-rprx-vision" {
		return fmt.Errorf("unsupported flow %q", cfg.Flow)
	}
	if cfg.Network != "" && cfg.Network != "tcp" && cfg.Network != "udp" {
		return fmt.Errorf("unsupported network %q", cfg.Network)
	}
	return nil
}

func (Adapter) Build(cfg config.OutboundConfig) (map[string]any, error) {
	out := map[string]any{
		"type":        "vless",
		"tag":         "proxy",
		"server":      cfg.Server,
		"server_port": cfg.Port,
		"uuid":        cfg.UUID,
	}
	if cfg.Flow != "" {
		out["flow"] = cfg.Flow
	}
	if cfg.Network != "" {
		out["network"] = cfg.Network
	}
	out["tls"] = buildTLS(cfg.TLS)

	return out, nil
}

func buildTLS(cfg *config.TLSConfig) map[string]any {
	out := map[string]any{
		"enabled": true,
	}
	if cfg == nil {
		return out
	}
	if cfg.ServerName != "" {
		out["server_name"] = cfg.ServerName
	}
	if cfg.Insecure {
		out["insecure"] = true
	}
	if cfg.Reality != nil && cfg.Reality.Enabled {
		out["reality"] = map[string]any{
			"enabled":    true,
			"public_key": cfg.Reality.PublicKey,
			"short_id":   cfg.Reality.ShortID,
		}
	}
	if cfg.Reality != nil && cfg.Reality.Fingerprint != "" {
		out["utls"] = map[string]any{
			"enabled":     true,
			"fingerprint": cfg.Reality.Fingerprint,
		}
	}
	return out
}
