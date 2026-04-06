# Ghostify Phase Status

## Implemented Foundation

- Phase 1: Sing-box execution proof
- Phase 2: runtime manager with explicit lifecycle states
- Phase 3: Ghostify-owned JSON config loader and validator
- Phase 4: Sing-box mapper layer
- Phase 5: end-to-end `ghostify run -c ...`
- Phase 6: protocol registry
- Phase 7: initial native pack with `socks5` and `vless`
- URI import foundation for `vless://`
- Added `vmess` protocol adapter and importer
- Added `trojan` protocol adapter and importer
- Phase 8 started: OpenVPN bridge runner and documentation
- SSH bridge runner added and documented
- UDP2RAW bridge runner added
- Bridge registry command exposes available bridges
- Phase 9 (observability) delivered via `health`, `status`, metrics documentation
- Phase 10 (packaging) kickoff with `scripts/build.sh` + packaging notes
- Bridge registry command exposes available bridges

## Deferred by Design

- Bridge adapters such as OpenVPN and SSH-family compatibility
- Experimental protocols and sidecar runtimes
- Persistent status storage and richer observability
- Packaging and release automation

## Current Supported Outbound Types

- `socks5`
- `vless`
