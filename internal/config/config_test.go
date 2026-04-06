package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadNormalizesDefaults(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "ghostify.json")
	data := `{
	  "outbound": {
	    "type": "socks5",
	    "server": "127.0.0.1",
	    "port": 1080
	  }
	}`
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if cfg.Inbound.Type != "mixed" {
		t.Fatalf("expected default inbound type, got %q", cfg.Inbound.Type)
	}
	if cfg.Inbound.Port != 1080 {
		t.Fatalf("expected default inbound port, got %d", cfg.Inbound.Port)
	}
	if cfg.LogLevel != "info" {
		t.Fatalf("expected default log level, got %q", cfg.LogLevel)
	}
}

func TestValidateRejectsUnsupportedOutbound(t *testing.T) {
	t.Parallel()

	cfg := Config{
		LogLevel: "info",
		Inbound: InboundConfig{
			Type:   "mixed",
			Listen: "127.0.0.1",
			Port:   1080,
		},
		Outbound: OutboundConfig{
			Type:   "openvpn",
			Server: "127.0.0.1",
			Port:   443,
		},
	}

	if err := cfg.Validate(); err == nil {
		t.Fatal("expected validation error, got nil")
	}
}

func TestValidateVLESSRequiresTLS(t *testing.T) {
	t.Parallel()

	cfg := Config{
		LogLevel: "info",
		Inbound: InboundConfig{
			Type:   "mixed",
			Listen: "127.0.0.1",
			Port:   1080,
		},
		Outbound: OutboundConfig{
			Type:   "vless",
			Server: "example.com",
			Port:   443,
			UUID:   "bf000d23-0752-40b4-affe-68f7707a9661",
		},
	}

	if err := cfg.Validate(); err == nil {
		t.Fatal("expected validation error, got nil")
	}
}
