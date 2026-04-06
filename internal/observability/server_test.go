package observability

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/quadraphony/ghostfy/internal/bridge"
)

func TestStatusEndpoint(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/status", nil)
	w := httptest.NewRecorder()

	srv := &Server{Status: func() map[string]any {
		return map[string]any{
			"protocols": []map[string]string{},
			"bridges":   []bridge.Metadata{},
		}
	}}
	srv.Handler().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var decoded map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &decoded); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if _, ok := decoded["metrics"]; !ok {
		t.Fatalf("metrics not present")
	}
}
