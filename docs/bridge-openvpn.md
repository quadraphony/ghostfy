# OpenVPN Bridge

Ghostify can run an OpenVPN sidecar as a trusted bridge.

## Config

```json
{
  "profile": "openvpn-bridge",
  "log_level": "info",
  "openvpn_config": "/etc/openvpn/client.conf",
  "args": ["--verb", "3"],
  "openvpn_path": "/usr/bin/openvpn"
}
```

## Run

```bash
ghostify openvpn-bridge -c examples/ghostify.openvpn.json
```

The command streams OpenVPN logs to stdout/stderr and stops cleanly on Ctrl+C.
