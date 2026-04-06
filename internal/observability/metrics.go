package observability

import (
	"sync"
	"time"
)

type Metrics struct {
	mu        sync.Mutex
	RunCount  int
	LastRun   time.Time
	LastState string
}

var Default = &Metrics{}

func (m *Metrics) Record(state string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.RunCount++
	m.LastRun = time.Now().UTC()
	m.LastState = state
}

func (m *Metrics) Snapshot() map[string]any {
	m.mu.Lock()
	defer m.mu.Unlock()
	return map[string]any{
		"run_count":  m.RunCount,
		"last_run":   m.LastRun.Format(time.RFC3339),
		"last_state": m.LastState,
	}
}
