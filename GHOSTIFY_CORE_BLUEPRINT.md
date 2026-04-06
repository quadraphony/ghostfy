# Ghostify-Core — Master Blueprint

## 1. Project Identity

**Name:** Ghostify-Core  
**Type:** CLI-first modular tunneling core  
**Primary Language:** Go  
**Frontend:** None for now  
**Execution Strategy:** Sing-box-backed first, Ghostify abstraction on top  
**Goal:** Build a clean, extensible tunneling engine that can eventually unify many proxy, VPN, transport, and stealth protocols under one Ghostify-controlled architecture.

---

## 2. Core Vision

Ghostify-Core is a **modular orchestration and extension layer** for tunneling and proxy technologies.

It is **not** meant to be a messy fork of Sing-box.

It is meant to:

- prove reliable runtime control over Sing-box first
- wrap execution with Ghostify lifecycle management
- introduce a Ghostify-native config system
- add a Ghostify adapter layer
- support many protocols safely over time
- keep architecture clean enough for future mobile and app embedding

---

## 3. Strategic Direction

Ghostify-Core must be built in this order:

### Stage A — Prove execution

Make sure Sing-box can be launched, monitored, and stopped correctly from Go.

### Stage B — Wrap runtime

Build Ghostify runtime lifecycle around Sing-box execution.

### Stage C — Add Ghostify config

Introduce a human-friendly config format that Ghostify owns.

### Stage D — Add mapping

Convert Ghostify config into Sing-box-compatible runtime config.

### Stage E — Add protocol modules

Introduce modular protocol abstractions and expand gradually.

### Stage F — Expand beyond direct native support

Where protocols are not directly available in Sing-box, support them through Ghostify bridge/adapters, sidecar strategies, importers, translators, compatibility wrappers, or experimental modules.

---

## 4. Core Non-Negotiable Rules

1. Do not modify Sing-box source code.
2. Do not fork Sing-box into a custom patch mess.
3. All Ghostify logic must live inside Ghostify-Core.
4. Every phase must compile and pass tests before the next phase starts.
5. No fake implementations pretending a protocol works.
6. No frontend or GUI in this stage.
7. Config must be human-friendly.
8. Validation must be strict.
9. Logging must be structured.
10. Runtime lifecycle must be explicit and testable.
11. Each protocol must have isolated validation logic.
12. Architecture must allow future growth without rewriting the foundation.
13. Protocol ambition is large, but first implementation scope must remain controlled.
14. Push to GitHub only after a stable tested phase.
15. Update docs after every stable phase.

---

## 5. High-Level Architecture

```text
Ghostify Config (ghostify.json)
        ↓
Config Loader
        ↓
Config Validator
        ↓
Ghostify Internal Model
        ↓
Protocol Registry
        ↓
Adapter / Mapper Layer
        ↓
Execution Backend
   ├── Sing-box native mapping
   ├── External bridge adapter
   ├── Sidecar runtime
   └── Experimental module
        ↓
Runtime Manager
        ↓
Logging / Status / Health / Stats
```

---

## 6. What Ghostify-Core Is Trying to Achieve

Ghostify-Core aims to become a **unified control plane** for many classes of tunneling and proxy technologies:

- classic proxy protocols
- censorship-resistance transports
- VPN-style tunnels
- stealth/obfuscation layers
- SSH-based tunneling strategies
- DNS tunneling strategies
- QUIC-based transports
- plugin/bridge-based compatibility modules
- future app embedding

This means the long-term goal is not "support one protocol."

The long-term goal is:

> build a clean engine that can absorb and control many protocols without becoming unmaintainable

---

## 7. Protocol Universe — Full Target Vision

This section lists the broad protocol universe Ghostify-Core should be designed to accommodate.

This is the target architecture vision, not a promise that all of them are implemented on day one.

---

### 7.1 Core Proxy / Standard Protocols

These are baseline support targets:

- SOCKS4
- SOCKS4A
- SOCKS5
- HTTP Proxy
- HTTPS Proxy
- HTTP CONNECT
- Mixed proxy modes

---

### 7.2 Modern Proxy Ecosystem Protocols

These are high-priority modern proxy protocols:

- Shadowsocks
- Shadowsocks 2022
- ShadowsocksR
- VMess
- VLESS
- Trojan
- NaiveProxy
- ShadowTLS
- AnyTLS

---

### 7.3 QUIC / UDP / Modern High-Performance Protocols

These are high-priority modern performance-focused protocols:

- Hysteria
- Hysteria2
- TUIC
- QUIC-based custom transport modes
- UDP-over-TCP wrapper strategies
- TCP Brutal-like transport control concepts
- multiplexing layers where appropriate

---

### 7.4 VPN / Tunnel Protocols

These are strategic tunnel targets:

- WireGuard
- OpenVPN
- TUN-based local VPN execution
- TProxy modes
- Redirect modes
- transparent proxy tunnel integration
- route-driven endpoint execution

---

### 7.5 SSH-Based Tunneling Family

These are highly valuable Ghostify expansion targets:

- SSH direct
- SSH proxy
- SSH reverse tunnel concepts
- SSH over TLS
- SSH over WebSocket
- SSH payload modes
- SSH-based custom encapsulation
- SSH + obfuscation combinations

---

### 7.6 Transport and Layering Options

Ghostify-Core should support protocol layering and transport variations:

- TCP
- UDP
- TLS
- mTLS
- WebSocket
- HTTP/2
- gRPC
- HTTP/3
- QUIC
- Reality-style security layering
- SNI-based routing strategies
- ALPN tuning
- domain fronting concepts where technically feasible
- CDN-friendly transport wrapping
- multiplexing
- fragmentation / packet behavior tuning where supported

---

### 7.7 Obfuscation / Stealth / Evasion Targets

These are valuable for long-term architecture support:

- Obfs4-style systems
- Meek-like strategies
- packet camouflage ideas
- TLS camouflage
- header obfuscation
- traffic shaping
- fingerprint tuning
- JA3 / TLS fingerprint control concepts
- stealth wrappers
- custom pluggable obfuscation modules

---

### 7.8 DNS / Slow Transport / Constrained-Network Targets

These are future expansion targets:

- DoH
- DoT
- DNS-over-QUIC
- DNS tunnel concepts
- SlowDNS
- DNSTT-style ideas
- constrained-bandwidth fallback strategies

---

### 7.9 Special / Advanced / Compatibility Targets

These are long-term compatibility or specialty targets:

- Tor integration
- WARP-style integration strategy
- Psiphon-style compatibility strategy
- UDP2RAW-style compatibility ideas
- KCP / mKCP-style transport concepts
- SMUX / YAMUX-style multiplex ideas
- Tailscale-related endpoint ideas where useful
- selector and balancing strategies
- URL test / health-based outbound selection
- failover chains
- multi-hop routing chains

---

## 8. Protocol Classification Model

Ghostify-Core must not treat all protocols the same.

Each protocol should be classified into one of these support classes:

### Class A — Native Sing-box-backed

Protocols Ghostify can directly map into Sing-box-supported runtime configuration.

### Class B — Native-adjacent

Protocols supported through Sing-box endpoint, inbound, outbound, or transport structures with Ghostify-specific wrapping.

### Class C — Bridge / Adapter

Protocols that are not cleanly native but can be supported through:

- config translators
- sidecar processes
- wrapper runtimes
- import/export compatibility
- external binary orchestration

### Class D — Experimental

Protocols or stealth systems that are research-heavy, unstable, or intentionally isolated until proven.

---

## 9. Initial Native-First Strategy

Ghostify-Core should begin with protocols that Sing-box already gives strong execution leverage for, because that reduces engineering uncertainty and proves the control-plane design faster.

The first architecture should assume that Ghostify can grow over the currently documented Sing-box-supported proxy/tunnel surface, then later expand into bridge-based support for non-native ecosystems.

---

## 10. Full Support Roadmap Model

Ghostify should track every protocol under one of these statuses:

- PLANNED
- ARCHITECTURE-READY
- NATIVE-MAPPABLE
- BRIDGE-READY
- EXPERIMENTAL
- IMPLEMENTED
- TESTED
- STABLE
- DEPRECATED
- REMOVED

---

## 11. First Build Scope vs Full Vision

### Full Vision

Support a huge protocol universe over time.

### First Build Scope

Do not implement all protocols immediately.

The first implementation exists to prove that the engine can control Sing-box cleanly.

That means Phase 1 should stay very small.

---

## 12. Phase Plan

---

# PHASE 1 — Sing-box Execution Proof

## Goal

Prove Sing-box works reliably under Ghostify control.

## Tasks

- prepare Sing-box binary strategy
- run Sing-box from Go using os/exec
- use a minimal valid test config
- capture stdout
- capture stderr
- detect startup failure
- detect crash
- stop process cleanly
- log everything properly

## Output

CLI command:

```bash
ghostify run-singbox-test
```

## Completion Criteria

- Sing-box can be launched from Ghostify
- output is captured
- stop flow works
- error flow works
- tests pass

---

# PHASE 2 — Ghostify Runtime Manager

## Goal

Add lifecycle control.

## Tasks

- runtime manager
- lifecycle states
- start
- stop
- restart
- graceful shutdown
- signal handling
- structured logs

## States

- INIT
- STARTING
- RUNNING
- STOPPING
- STOPPED
- ERROR

---

# PHASE 3 — Ghostify Config System

## Goal

Introduce Ghostify-owned config.

## Tasks

- define JSON schema
- implement loader
- implement validator
- fail fast on invalid fields
- normalize data

## Example

```json
{
  "profile": "test",
  "log_level": "info",
  "outbound": {
    "type": "socks5",
    "server": "127.0.0.1",
    "port": 1080
  }
}
```

---

# PHASE 4 — Mapper Layer

## Goal

Convert Ghostify config into Sing-box config.

## Tasks

- internal config model
- mapper package
- conversion logic
- validation of supported combinations
- generated config output
- mapper tests

---

# PHASE 5 — First End-to-End Run

## Goal

Run full Ghostify flow.

## Flow

```text
ghostify run -c ghostify.json
→ load config
→ validate
→ map to sing-box config
→ launch sing-box
→ stream logs
→ stop gracefully
```

---

# PHASE 6 — Protocol Registry

## Goal

Introduce modular protocol abstraction.

## Tasks

- protocol interface
- protocol registry
- protocol-specific validators
- separation of protocol concerns
- preparation for multi-protocol expansion

## Suggested Interface

```go
type ProtocolAdapter interface {
    Name() string
    Validate(cfg any) error
    Build(cfg any) (any, error)
}
```

---

# PHASE 7 — Controlled Native Expansion

## Goal

Expand the first stable protocol pack.

## Recommended Order

1. SOCKS5
2. HTTP/HTTPS proxy
3. Shadowsocks
4. VMess
5. VLESS
6. Trojan
7. Naive
8. ShadowTLS
9. Hysteria2
10. TUIC
11. AnyTLS
12. WireGuard-oriented support model
13. SSH outbound-oriented support model
14. Tor / DNS / selector / urltest control models

---

# PHASE 8 — Bridge / Compatibility Layer

## Goal

Prepare support for protocols outside direct native mapping.

## Target Ideas

- OpenVPN compatibility
- ShadowsocksR compatibility
- SSH variants
- SlowDNS / DNSTT concepts
- WARP strategy
- Psiphon-like strategy
- UDP2RAW ideas
- KCP / mKCP concepts
- stealth wrappers

## Methods

- sidecar process model
- import/export translator
- compatibility wrapper
- external adapter interface
- experimental module isolation

---

# PHASE 9 — Observability

## Goal

Improve operator usability.

## Tasks

- structured logs
- health command
- status command
- runtime summary
- basic stats
- crash reason surfaces
- verbose debug mode

---

# PHASE 10 — Packaging

## Goal

Distribute Ghostify-Core cleanly.

## Tasks

- Linux builds
- Windows builds
- version metadata
- build scripts
- release docs
- install docs

---

## 13. Protocol Registry Design Requirement

Each protocol must define:

- name
- category
- status
- support_class
- validator
- mapper
- capability flags
- test coverage
- docs

Suggested metadata shape:

```go
type ProtocolMetadata struct {
    Name         string
    Category     string
    Status       string
    SupportClass string
    Experimental bool
}
```

---

## 14. Suggested Protocol Categories Table

### A. Proxy Basics

- socks4
- socks4a
- socks5
- http
- https
- mixed

### B. Modern Proxy Core

- shadowsocks
- shadowsocks2022
- shadowsocksr
- vmess
- vless
- trojan
- naive
- shadowtls
- anytls

### C. QUIC / Performance

- hysteria
- hysteria2
- tuic
- quic-custom

### D. Tunnel / VPN

- wireguard
- openvpn
- tun
- tproxy
- redirect

### E. SSH Family

- ssh-direct
- ssh-proxy
- ssh-tls
- ssh-ws
- ssh-payload

### F. DNS / Slow Links

- doh
- dot
- doq
- dns-tunnel
- slowdns
- dnstt

### G. Stealth / Obfuscation

- obfs4
- meek
- tls-camouflage
- fingerprint-tuned
- fronted
- custom-obfs

### H. Advanced Routing / Control

- tor
- warp
- psiphon
- selector
- urltest
- chain
- failover
- balancer

### I. Experimental / Compatibility

- udp2raw
- kcp
- mkcp
- smux
- yamux
- custom-sidecar
- external-adapter

---

## 15. Repo Structure

```text
ghostify-core/
├── cmd/
│   └── ghostify/
│       └── main.go
├── internal/
│   ├── app/
│   ├── config/
│   ├── runtime/
│   ├── logging/
│   ├── stats/
│   ├── adapters/
│   │   ├── singbox/
│   │   ├── bridge/
│   │   └── external/
│   ├── protocols/
│   │   ├── registry/
│   │   ├── common/
│   │   ├── socks5/
│   │   ├── http/
│   │   ├── shadowsocks/
│   │   ├── vmess/
│   │   ├── vless/
│   │   ├── trojan/
│   │   ├── naive/
│   │   ├── shadowtls/
│   │   ├── hysteria2/
│   │   ├── tuic/
│   │   ├── wireguard/
│   │   ├── ssh/
│   │   ├── tor/
│   │   ├── openvpn/
│   │   ├── slowdns/
│   │   ├── dnstt/
│   │   ├── obfs/
│   │   └── experimental/
│   └── transports/
│       ├── tcp/
│       ├── udp/
│       ├── tls/
│       ├── ws/
│       ├── grpc/
│       ├── quic/
│       └── http2/
├── examples/
├── docs/
├── test/
├── scripts/
├── go.mod
└── README.md
```

---

## 16. Definition of Done

A phase is done only when:

- code compiles
- tests pass
- logs are readable
- docs are updated
- no fake support claims exist
- phase boundary is respected

---

## 17. What Must Not Happen

- do not implement 20 protocols in Phase 1
- do not create fake protocol support
- do not mix runtime and protocol logic into one file
- do not hardcode Sing-box assumptions everywhere
- do not create a god-manager
- do not skip tests
- do not claim compatibility without proof
- do not let large protocol ambition destroy architecture quality

---

## 18. First Real Milestone

The first real milestone is not "support everything."

The first real milestone is:

> Ghostify-Core can launch Sing-box from Go, control lifecycle, validate Ghostify config, map one supported protocol correctly, stream logs, and stop cleanly.

That is the foundation.

---

## 19. Immediate Instruction to AI Dev

Start with **Phase 1 only**.

Do not jump ahead.

Do not add protocol overload early.

First prove Sing-box execution and Ghostify lifecycle control.
