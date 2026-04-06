package importers

import (
	"encoding/base64"
	"encoding/json"
	"testing"
)

func TestImportVLESSRealityURI(t *testing.T) {
	t.Parallel()

	raw := "vless://43488128-319e-f480-64ea-0acdc712e2a8@51.68.155.153:443?encryption=none&flow=xtls-rprx-vision&fp=chrome&pbk=_xzl59bcUYD9QUSKFyboqCC_9eUnaXUUKWA19oQWFHU&security=reality&sid=8d17b7ebab7b99d7&sni=www.allegro.pl&type=tcp#PL-vless"
	cfg, err := Import(raw)
	if err != nil {
		t.Fatalf("Import returned error: %v", err)
	}

	if cfg.Outbound.Type != "vless" {
		t.Fatalf("expected vless, got %q", cfg.Outbound.Type)
	}
	if cfg.Outbound.Server != "51.68.155.153" {
		t.Fatalf("unexpected server: %q", cfg.Outbound.Server)
	}
	if cfg.Outbound.TLS == nil || cfg.Outbound.TLS.Reality == nil {
		t.Fatal("expected reality tls config")
	}
	if cfg.Outbound.TLS.Reality.PublicKey == "" {
		t.Fatal("expected reality public key")
	}
}

func TestImportVMESSURI(t *testing.T) {
	t.Parallel()

	payload := map[string]string{
		"v":    "2",
		"ps":   "vmess",
		"add":  "example.com",
		"port": "443",
		"id":   "uuid",
		"aid":  "0",
		"scy":  "auto",
		"net":  "tcp",
		"type": "tcp",
		"tls":  "tls",
		"sni":  "example.com",
	}

	raw, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}

	uri := "vmess://" + base64.StdEncoding.EncodeToString(raw)
	cfg, err := Import(uri)
	if err != nil {
		t.Fatalf("Import returned error: %v", err)
	}
	if cfg.Outbound.Type != "vmess" {
		t.Fatalf("expected vmess, got %q", cfg.Outbound.Type)
	}
	if cfg.Outbound.Security != "auto" {
		t.Fatalf("expected security auto, got %q", cfg.Outbound.Security)
	}
}

func TestImportTrojanURI(t *testing.T) {
	t.Parallel()

	raw := "trojan://password@51.68.155.153:443?security=tls&fingerprint=chrome&sni=www.allegro.pl#trojan"
	cfg, err := Import(raw)
	if err != nil {
		t.Fatalf("Import returned error: %v", err)
	}

	if cfg.Outbound.Type != "trojan" {
		t.Fatalf("expected trojan, got %q", cfg.Outbound.Type)
	}
	if cfg.Outbound.Password != "password" {
		t.Fatalf("unexpected password: %q", cfg.Outbound.Password)
	}
	if cfg.Outbound.TLS == nil || !cfg.Outbound.TLS.Enabled {
		t.Fatal("expected TLS to be enabled")
	}
}

func TestImportRejectsUnsupportedScheme(t *testing.T) {
	t.Parallel()

	if _, err := Import("openvpn://example"); err == nil {
		t.Fatal("expected error for unsupported scheme")
	}
}

func TestImportRejectsUnsupportedVLESSSecurity(t *testing.T) {
	t.Parallel()

	raw := "vless://id@example.com:443?encryption=none&security=ws"
	if _, err := Import(raw); err == nil {
		t.Fatal("expected security error")
	}
}
