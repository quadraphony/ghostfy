package app

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestListProtocols(t *testing.T) {
	t.Parallel()

	app := New()
	protocols := app.ListProtocols()
	if len(protocols) < 2 {
		t.Fatalf("expected protocol list, got %d", len(protocols))
	}
}

func TestWriteRuntimeConfig(t *testing.T) {
	t.Parallel()

	path, err := writeRuntimeConfig(map[string]any{"log": map[string]any{"level": "info"}})
	if err != nil {
		t.Fatalf("writeRuntimeConfig returned error: %v", err)
	}
	defer os.Remove(path)

	if filepath.Ext(path) != ".json" {
		t.Fatalf("expected json file, got %s", path)
	}
}

func TestImportURI(t *testing.T) {
	t.Parallel()

	app := New()
	var out bytes.Buffer
	raw := "vless://43488128-319e-f480-64ea-0acdc712e2a8@51.68.155.153:443?encryption=none&flow=xtls-rprx-vision&fp=chrome&pbk=_xzl59bcUYD9QUSKFyboqCC_9eUnaXUUKWA19oQWFHU&security=reality&sid=8d17b7ebab7b99d7&sni=www.allegro.pl&type=tcp#PL-vless"
	if err := app.ImportURI(raw, &out); err != nil {
		t.Fatalf("ImportURI returned error: %v", err)
	}
	if out.Len() == 0 {
		t.Fatal("expected importer output")
	}
}
