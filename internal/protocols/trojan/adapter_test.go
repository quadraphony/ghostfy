package trojan

import (
	"testing"

	"github.com/quadraphony/ghostfy/internal/config"
)

func TestBuildTrojan(t *testing.T) {
	t.Parallel()

	adapter := Adapter{}

	cfg := config.OutboundConfig{
		Type:     "trojan",
		Server:   "51.68.155.153",
		Port:     443,
		Password: "password",
		TLS: &config.TLSConfig{
			Enabled:    true,
			ServerName: "www.allegro.pl",
		},
	}

	out, err := adapter.Build(cfg)
	if err != nil {
		t.Fatalf("build returned error: %v", err)
	}
	if out["type"] != "trojan" {
		t.Fatalf("expected trojan, got %v", out["type"])
	}
}
