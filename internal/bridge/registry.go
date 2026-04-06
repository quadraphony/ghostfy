package bridge

import "sort"

type Metadata struct {
	Name     string
	Category string
	Status   string
}

type Registry struct {
	entries map[string]Metadata
}

func New(entries ...Metadata) *Registry {
	r := &Registry{entries: make(map[string]Metadata, len(entries))}
	for _, entry := range entries {
		r.entries[entry.Name] = entry
	}
	return r
}

func Default() *Registry {
	return New(
		Metadata{Name: "openvpn-bridge", Category: "vpn", Status: "Stable"},
		Metadata{Name: "ssh-bridge", Category: "tunnel", Status: "Stable"},
		Metadata{Name: "udp2raw-bridge", Category: "stealth", Status: "Experimental"},
	)
}

func (r *Registry) List() []Metadata {
	list := make([]Metadata, 0, len(r.entries))
	for _, entry := range r.entries {
		list = append(list, entry)
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Name < list[j].Name
	})
	return list
}
