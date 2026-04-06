package bridge

import "testing"

func TestRegistryList(t *testing.T) {
	t.Parallel()

	reg := New(
		Metadata{Name: "openvpn", Category: "vpn", Status: "Ready"},
		Metadata{Name: "ssh", Category: "tunnel", Status: "Ready"},
	)

	if len(reg.List()) != 2 {
		t.Fatalf("expected two entries, got %d", len(reg.List()))
	}
}
