package singbox

import (
	"testing"

	"github.com/quadraphony/ghostfy/internal/config"
)

func TestBuildSocksConfig(t *testing.T) {
	t.Parallel()

	mapper := NewMapper()
	cfg := config.Config{
		LogLevel: "info",
		Inbound: config.InboundConfig{
			Type:   "mixed",
			Listen: "127.0.0.1",
			Port:   1080,
		},
		Outbound: config.OutboundConfig{
			Type:   "socks5",
			Server: "127.0.0.1",
			Port:   2080,
		},
	}

	built, err := mapper.Build(cfg)
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}

	outbounds := built["outbounds"].([]map[string]any)
	if outbounds[0]["type"] != "socks" {
		t.Fatalf("expected socks outbound, got %v", outbounds[0]["type"])
	}
}

func TestBuildVLESSConfig(t *testing.T) {
	t.Parallel()

	mapper := NewMapper()
	cfg := config.Config{
		LogLevel: "info",
		Inbound: config.InboundConfig{
			Type:   "mixed",
			Listen: "127.0.0.1",
			Port:   1080,
		},
		Outbound: config.OutboundConfig{
			Type:    "vless",
			Server:  "example.com",
			Port:    443,
			UUID:    "bf000d23-0752-40b4-affe-68f7707a9661",
			Flow:    "xtls-rprx-vision",
			Network: "tcp",
			TLS: &config.TLSConfig{
				Enabled:    true,
				ServerName: "example.com",
				Reality: &config.RealityConfig{
					Enabled:     true,
					PublicKey:   "pubkey",
					ShortID:     "abcd",
					Fingerprint: "chrome",
				},
			},
			Transport: &config.Transport{
				Type: "tcp",
			},
		},
	}

	built, err := mapper.Build(cfg)
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}

	outbounds := built["outbounds"].([]map[string]any)
	if outbounds[0]["type"] != "vless" {
		t.Fatalf("expected vless outbound, got %v", outbounds[0]["type"])
	}
}
