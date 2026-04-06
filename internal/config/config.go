package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Profile  string         `json:"profile"`
	LogLevel string         `json:"log_level"`
	Inbound  InboundConfig  `json:"inbound"`
	Outbound OutboundConfig `json:"outbound"`
}

type InboundConfig struct {
	Type     string `json:"type"`
	Listen   string `json:"listen"`
	Port     int    `json:"port"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type OutboundConfig struct {
	Type      string     `json:"type"`
	Server    string     `json:"server"`
	Port      int        `json:"port"`
	Username  string     `json:"username,omitempty"`
	Password  string     `json:"password,omitempty"`
	Version   string     `json:"version,omitempty"`
	UUID      string     `json:"uuid,omitempty"`
	Flow      string     `json:"flow,omitempty"`
	Network   string     `json:"network,omitempty"`
	AlterID   int        `json:"alter_id,omitempty"`
	Security  string     `json:"security,omitempty"`
	TLS       *TLSConfig `json:"tls,omitempty"`
	Transport *Transport `json:"transport,omitempty"`
}

type TLSConfig struct {
	Enabled    bool           `json:"enabled"`
	ServerName string         `json:"server_name,omitempty"`
	Insecure   bool           `json:"insecure,omitempty"`
	Reality    *RealityConfig `json:"reality,omitempty"`
}

type RealityConfig struct {
	Enabled     bool   `json:"enabled"`
	PublicKey   string `json:"public_key"`
	ShortID     string `json:"short_id,omitempty"`
	Fingerprint string `json:"fingerprint,omitempty"`
}

type Transport struct {
	Type string `json:"type"`
	Path string `json:"path,omitempty"`
	Host string `json:"host,omitempty"`
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

	cfg.Normalize()
	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (c *Config) Normalize() {
	if c.Profile == "" {
		c.Profile = "default"
	}
	if c.LogLevel == "" {
		c.LogLevel = "info"
	}
	if c.Inbound.Type == "" {
		c.Inbound.Type = "mixed"
	}
	if c.Inbound.Listen == "" {
		c.Inbound.Listen = "127.0.0.1"
	}
	if c.Inbound.Port == 0 {
		c.Inbound.Port = 1080
	}
	c.Inbound.Type = strings.ToLower(c.Inbound.Type)
	c.Outbound.Type = strings.ToLower(c.Outbound.Type)
	c.Outbound.Network = strings.ToLower(c.Outbound.Network)
	c.Outbound.Version = strings.ToLower(c.Outbound.Version)
	c.Outbound.Security = strings.ToLower(c.Outbound.Security)
	if c.Outbound.Type == "vmess" && c.Outbound.Security == "" {
		c.Outbound.Security = "auto"
	}
	if c.Outbound.TLS != nil && c.Outbound.TLS.Reality != nil {
		c.Outbound.TLS.Reality.Fingerprint = strings.ToLower(c.Outbound.TLS.Reality.Fingerprint)
	}
	if c.TransportOrNil() != nil {
		c.Outbound.Transport.Type = strings.ToLower(c.Outbound.Transport.Type)
	}
}

func (c Config) Validate() error {
	if c.LogLevel == "" {
		return errors.New("log_level is required")
	}
	switch c.LogLevel {
	case "trace", "debug", "info", "warn", "error":
	default:
		return fmt.Errorf("unsupported log_level %q", c.LogLevel)
	}

	if c.Inbound.Port < 1 || c.Inbound.Port > 65535 {
		return fmt.Errorf("inbound.port must be between 1 and 65535")
	}

	switch c.Inbound.Type {
	case "mixed", "socks", "http":
	default:
		return fmt.Errorf("unsupported inbound.type %q", c.Inbound.Type)
	}

	if c.Outbound.Server == "" {
		return errors.New("outbound.server is required")
	}
	if c.Outbound.Port < 1 || c.Outbound.Port > 65535 {
		return fmt.Errorf("outbound.port must be between 1 and 65535")
	}

	switch c.Outbound.Type {
	case "vmess":
		if c.Outbound.UUID == "" {
			return errors.New("outbound.uuid is required for vmess")
		}
		if c.Outbound.Port < 1 || c.Outbound.Port > 65535 {
			return fmt.Errorf("outbound.port must be between 1 and 65535")
		}
		if c.Outbound.Security != "" {
			switch c.Outbound.Security {
			case "auto", "none", "zero", "aes-128-gcm", "chacha20-poly1305", "aes-128-ctr":
			default:
				return fmt.Errorf("unsupported vmess security %q", c.Outbound.Security)
			}
		}
		switch c.Outbound.Network {
		case "", "tcp", "ws", "h2":
		default:
			return fmt.Errorf("unsupported vmess network %q", c.Outbound.Network)
		}
		if c.Outbound.Transport != nil && c.Outbound.Transport.Type != "" {
			switch c.Outbound.Transport.Type {
			case "tcp", "ws", "http", "h2":
			default:
				return fmt.Errorf("unsupported vmess transport.type %q", c.Outbound.Transport.Type)
			}
		}
		return nil
	case "socks5":
		if c.Outbound.Version == "" {
			return nil
		}
		switch c.Outbound.Version {
		case "4", "4a", "5":
			return nil
		default:
			return fmt.Errorf("unsupported outbound.version %q", c.Outbound.Version)
		}
	case "vless":
		if c.Outbound.UUID == "" {
			return errors.New("outbound.uuid is required for vless")
		}
		if c.Outbound.Flow != "" && c.Outbound.Flow != "xtls-rprx-vision" {
			return fmt.Errorf("unsupported outbound.flow %q", c.Outbound.Flow)
		}
		if c.Outbound.Network != "" && c.Outbound.Network != "tcp" && c.Outbound.Network != "udp" {
			return fmt.Errorf("unsupported outbound.network %q", c.Outbound.Network)
		}
		if c.Outbound.TLS == nil || !c.Outbound.TLS.Enabled {
			return errors.New("outbound.tls.enabled must be true for vless")
		}
		if c.Outbound.TLS.Reality != nil && c.Outbound.TLS.Reality.Enabled && c.Outbound.TLS.Reality.PublicKey == "" {
			return errors.New("outbound.tls.reality.public_key is required for reality")
		}
		if c.Outbound.Transport != nil && c.Outbound.Transport.Type != "" {
			return fmt.Errorf("unsupported outbound.transport.type %q", c.Outbound.Transport.Type)
		}
	case "trojan":
		if c.Outbound.Password == "" {
			return errors.New("outbound.password is required for trojan")
		}
		if c.Outbound.Security == "" && (c.Outbound.TLS == nil || !c.Outbound.TLS.Enabled) {
			return errors.New("tls.enabled must be true for secure trojan")
		}
		return nil
	case "trusttunnel":
		if c.Outbound.Server == "" {
			return errors.New("outbound.server is required for trusttunnel")
		}
		if c.Outbound.Port < 1 || c.Outbound.Port > 65535 {
			return fmt.Errorf("outbound.port must be between 1 and 65535")
		}
		return nil
	default:
		return fmt.Errorf("unsupported outbound.type %q", c.Outbound.Type)
	}

	return nil
}

func (c Config) TransportOrNil() *Transport {
	return c.Outbound.Transport
}
