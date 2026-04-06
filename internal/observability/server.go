package observability

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
)

type StatusProvider func() map[string]any

type Server struct {
	Addr   string
	Status StatusProvider
}

func (s *Server) ListenAndServe() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		payload := map[string]any{}
		if s.Status != nil {
			payload = s.Status()
		}
		payload["metrics"] = Default.Snapshot()
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(payload)
	})

	listener, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return fmt.Errorf("listen %s: %w", s.Addr, err)
	}
	return http.Serve(listener, mux)
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		payload := map[string]any{}
		if s.Status != nil {
			payload = s.Status()
		}
		payload["metrics"] = Default.Snapshot()
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(payload)
	})
	return mux
}
