package observability

import (
	"testing"
)

func TestMetricsSnapshot(t *testing.T) {
	t.Parallel()

	Default.Record("test")
	snapshot := Default.Snapshot()
	if snapshot["last_state"] != "test" {
		t.Fatalf("expected state test, got %v", snapshot["last_state"])
	}
	if snapshot["run_count"].(int) == 0 {
		t.Fatalf("expected run_count > 0")
	}
}
