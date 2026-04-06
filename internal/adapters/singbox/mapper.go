package singbox

import (
	"fmt"

	"github.com/quadraphony/ghostfy/internal/config"
	"github.com/quadraphony/ghostfy/internal/protocols/registry"
	"github.com/quadraphony/ghostfy/internal/protocols/socks5"
	"github.com/quadraphony/ghostfy/internal/protocols/trojan"
	"github.com/quadraphony/ghostfy/internal/protocols/vless"
	"github.com/quadraphony/ghostfy/internal/protocols/vmess"
)

type Mapper struct {
	registry *registry.Registry
}

func NewMapper() *Mapper {
	return &Mapper{
		registry: registry.New(
			socks5.Adapter{},
			vless.Adapter{},
			vmess.Adapter{},
			trojan.Adapter{},
		),
	}
}

func (m *Mapper) Registry() *registry.Registry {
	return m.registry
}

func (m *Mapper) Build(cfg config.Config) (map[string]any, error) {
	adapter, err := m.registry.Get(cfg.Outbound.Type)
	if err != nil {
		return nil, err
	}
	if err := adapter.Validate(cfg.Outbound); err != nil {
		return nil, fmt.Errorf("validate %s outbound: %w", cfg.Outbound.Type, err)
	}

	outbound, err := adapter.Build(cfg.Outbound)
	if err != nil {
		return nil, fmt.Errorf("build %s outbound: %w", cfg.Outbound.Type, err)
	}

	return map[string]any{
		"log": map[string]any{
			"level":     cfg.LogLevel,
			"timestamp": true,
		},
		"inbounds": []map[string]any{
			buildInbound(cfg.Inbound),
		},
		"outbounds": []map[string]any{
			outbound,
			{
				"type": "direct",
				"tag":  "direct",
			},
		},
		"route": map[string]any{
			"final": "proxy",
		},
	}, nil
}

func buildInbound(cfg config.InboundConfig) map[string]any {
	inbound := map[string]any{
		"type":        cfg.Type,
		"tag":         "local-in",
		"listen":      cfg.Listen,
		"listen_port": cfg.Port,
	}

	if cfg.Username != "" {
		inbound["users"] = []map[string]any{
			{
				"username": cfg.Username,
				"password": cfg.Password,
			},
		}
	}

	return inbound
}
