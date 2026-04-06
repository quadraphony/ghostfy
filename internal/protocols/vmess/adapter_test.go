package vmess

import (
	"testing"

	"github.com/quadraphony/ghostfy/internal/config"
)

func TestBuildBasic(t *testing.T) {
	t.Parallel()

	adapter := Adapter{}

	cfg := config.OutboundConfig{
		Type:     "vmess",
		Server:   "example.com",
		Port:     443,
		UUID:     "id",
		Security: "auto",
	}

	out, err := adapter.Build(cfg)
	if err != nil {
		t.Fatalf("build returned error: %v", err)
	}

	if out["type"] != "vmess" {
		t.Fatalf("expected vmess export, got %v", out["type"])
	}
}
