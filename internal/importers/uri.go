package importers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"

	"github.com/quadraphony/ghostfy/internal/config"
)

func Import(raw string) (config.Config, error) {
	parsed, err := url.Parse(strings.TrimSpace(raw))
	if err != nil {
		return config.Config{}, fmt.Errorf("parse uri: %w", err)
	}

	switch strings.ToLower(parsed.Scheme) {
	case "vless":
		return importVLESS(parsed)
	case "vmess":
		return importVMess(parsed)
	case "trojan":
		return importTrojan(parsed)
	default:
		return config.Config{}, fmt.Errorf("unsupported uri scheme %q", parsed.Scheme)
	}
}

type vmessPayload struct {
	Version  string `json:"v"`
	Remark   string `json:"ps"`
	Address  string `json:"add"`
	Port     string `json:"port"`
	ID       string `json:"id"`
	AlterID  string `json:"aid"`
	Security string `json:"scy"`
	Network  string `json:"net"`
	Type     string `json:"type"`
	Host     string `json:"host"`
	Path     string `json:"path"`
	TLS      string `json:"tls"`
	Sni      string `json:"sni"`
}

func importVMess(parsed *url.URL) (config.Config, error) {
	payload := parsed.Opaque
	if payload == "" {
		payload = parsed.Host
	}
	if payload == "" {
		return config.Config{}, fmt.Errorf("vmess uri missing payload")
	}

	data, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		return config.Config{}, fmt.Errorf("decode payload: %w", err)
	}

	var vm vmessPayload
	if err := json.Unmarshal(data, &vm); err != nil {
		return config.Config{}, fmt.Errorf("decode payload json: %w", err)
	}

	port, err := strconv.Atoi(vm.Port)
	if err != nil {
		return config.Config{}, fmt.Errorf("invalid port %q", vm.Port)
	}

	alterID := 0
	if vm.AlterID != "" {
		alterID, err = strconv.Atoi(vm.AlterID)
		if err != nil {
			return config.Config{}, fmt.Errorf("invalid alter_id %q", vm.AlterID)
		}
	}

	cfg := config.Config{
		Profile:  profileFromFragment(parsed.Fragment),
		LogLevel: "info",
		Inbound: config.InboundConfig{
			Type:   "mixed",
			Listen: "127.0.0.1",
			Port:   1080,
		},
		Outbound: config.OutboundConfig{
			Type:     "vmess",
			Server:   vm.Address,
			Port:     port,
			UUID:     vm.ID,
			AlterID:  alterID,
			Security: strings.ToLower(vm.Security),
			Network:  strings.ToLower(vm.Network),
			Transport: &config.Transport{
				Type: strings.ToLower(vm.Type),
				Path: vm.Path,
				Host: vm.Host,
			},
			TLS: &config.TLSConfig{
				Enabled:    strings.EqualFold(vm.TLS, "tls"),
				ServerName: vm.Sni,
			},
		},
	}

	cfg.Normalize()
	if err := cfg.Validate(); err != nil {
		return config.Config{}, err
	}

	return cfg, nil
}

func importVLESS(parsed *url.URL) (config.Config, error) {
	host := parsed.Hostname()
	if host == "" {
		return config.Config{}, fmt.Errorf("vless uri missing host")
	}

	port, err := parsePort(parsed)
	if err != nil {
		return config.Config{}, err
	}

	uuid := ""
	if parsed.User != nil {
		uuid = parsed.User.Username()
	}
	if uuid == "" {
		return config.Config{}, fmt.Errorf("vless uri missing uuid")
	}

	query := parsed.Query()
	security := strings.ToLower(query.Get("security"))
	fp := strings.ToLower(query.Get("fp"))
	transportType := strings.ToLower(query.Get("type"))
	if transportType == "" {
		transportType = "tcp"
	}

	cfg := config.Config{
		Profile:  profileFromFragment(parsed.Fragment),
		LogLevel: "info",
		Inbound: config.InboundConfig{
			Type:   "mixed",
			Listen: "127.0.0.1",
			Port:   1080,
		},
		Outbound: config.OutboundConfig{
			Type:    "vless",
			Server:  host,
			Port:    port,
			UUID:    uuid,
			Flow:    query.Get("flow"),
			Network: transportType,
		},
	}

	switch security {
	case "", "tls":
		cfg.Outbound.TLS = &config.TLSConfig{
			Enabled:    true,
			ServerName: query.Get("sni"),
		}
	case "reality":
		cfg.Outbound.TLS = &config.TLSConfig{
			Enabled:    true,
			ServerName: query.Get("sni"),
			Reality: &config.RealityConfig{
				Enabled:     true,
				PublicKey:   query.Get("pbk"),
				ShortID:     query.Get("sid"),
				Fingerprint: fp,
			},
		}
	default:
		return config.Config{}, fmt.Errorf("unsupported vless security %q", security)
	}

	if query.Get("encryption") != "" && query.Get("encryption") != "none" {
		return config.Config{}, fmt.Errorf("unsupported vless encryption %q", query.Get("encryption"))
	}

	cfg.Normalize()
	if err := cfg.Validate(); err != nil {
		return config.Config{}, err
	}

	return cfg, nil
}

func importTrojan(parsed *url.URL) (config.Config, error) {
	password := ""
	if parsed.User != nil {
		if user := parsed.User.Username(); user != "" {
			password = user
		} else {
			password, _ = parsed.User.Password()
		}
	}
	if password == "" {
		return config.Config{}, fmt.Errorf("trojan uri missing password")
	}

	host := parsed.Hostname()
	if host == "" {
		return config.Config{}, fmt.Errorf("trojan uri missing host")
	}

	port, err := parsePort(parsed)
	if err != nil {
		return config.Config{}, err
	}

	query := parsed.Query()

	cfg := config.Config{
		Profile:  profileFromFragment(parsed.Fragment),
		LogLevel: "info",
		Inbound: config.InboundConfig{
			Type:   "mixed",
			Listen: "127.0.0.1",
			Port:   1080,
		},
		Outbound: config.OutboundConfig{
			Type:     "trojan",
			Server:   host,
			Port:     port,
			Password: password,
			TLS: &config.TLSConfig{
				Enabled:    true,
				ServerName: query.Get("sni"),
			},
		},
	}

	cfg.Normalize()
	if err := cfg.Validate(); err != nil {
		return config.Config{}, err
	}

	return cfg, nil
}

func parsePort(parsed *url.URL) (int, error) {
	portText := parsed.Port()
	if portText == "" {
		return 0, fmt.Errorf("uri missing port")
	}

	port, err := strconv.Atoi(portText)
	if err != nil {
		return 0, fmt.Errorf("invalid port %q", portText)
	}
	if port < 1 || port > 65535 {
		return 0, fmt.Errorf("port must be between 1 and 65535")
	}
	return port, nil
}

func profileFromFragment(fragment string) string {
	fragment = strings.TrimSpace(fragment)
	if fragment == "" {
		return "imported"
	}

	var b strings.Builder
	for _, r := range fragment {
		switch {
		case r >= 'a' && r <= 'z':
			b.WriteRune(r)
		case r >= 'A' && r <= 'Z':
			b.WriteRune(r + ('a' - 'A'))
		case r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == '-' || r == '_':
			b.WriteRune(r)
		case r == ' ':
			b.WriteRune('-')
		}
	}

	profile := strings.Trim(b.String(), "-_")
	if profile == "" {
		host, _, err := net.SplitHostPort("placeholder:0")
		if err == nil && host != "" {
			return host
		}
		return "imported"
	}
	return profile
}
