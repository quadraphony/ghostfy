package trusttunnel

import (
	"testing"

	"github.com/quadraphony/ghostfy/internal/config"
)

func TestTrustTunnelBuild(t *testing.T) {
	t.Parallel()

	adapter := Adapter{}
	cfg := config.OutboundConfig{
		Type:     "trusttunnel",
		Server:   "51.68.155.153",
		Port:     443,
		Password: "password",
		Flow:     "relay",
	}

	out, err := adapter.Build(cfg)
	if err != nil {
		t.Fatalf("build error: %v", err)
	}
	if out["type"] != "trusttunnel" {
		t.Fatalf("unexpected type %v", out["type"])
	}
}
