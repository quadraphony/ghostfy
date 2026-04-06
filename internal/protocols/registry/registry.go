package registry

import (
	"fmt"
	"sort"

	"github.com/quadraphony/ghostfy/internal/config"
)

type ProtocolMetadata struct {
	Name         string
	Category     string
	Status       string
	SupportClass string
	Experimental bool
}

type ProtocolAdapter interface {
	Name() string
	Metadata() ProtocolMetadata
	Validate(cfg config.OutboundConfig) error
	Build(cfg config.OutboundConfig) (map[string]any, error)
}

type Registry struct {
	adapters map[string]ProtocolAdapter
}

func New(adapters ...ProtocolAdapter) *Registry {
	r := &Registry{adapters: make(map[string]ProtocolAdapter, len(adapters))}
	for _, adapter := range adapters {
		r.adapters[adapter.Name()] = adapter
	}
	return r
}

func (r *Registry) Get(name string) (ProtocolAdapter, error) {
	adapter, ok := r.adapters[name]
	if !ok {
		return nil, fmt.Errorf("protocol %q is not registered", name)
	}

	return adapter, nil
}

func (r *Registry) List() []ProtocolMetadata {
	names := make([]string, 0, len(r.adapters))
	for name := range r.adapters {
		names = append(names, name)
	}
	sort.Strings(names)

	out := make([]ProtocolMetadata, 0, len(names))
	for _, name := range names {
		out = append(out, r.adapters[name].Metadata())
	}
	return out
}
