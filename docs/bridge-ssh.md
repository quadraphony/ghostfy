# SSH Bridge

Ghostify can start an SSH client as a bridge for tunneling.

## Config

```json
{
  "profile": "ssh-bridge",
  "log_level": "info",
  "connection": "user@host -p 22",
  "args": ["-N", "-L", "1080:127.0.0.1:1080"]
}
```

## Run

```bash
ghostify ssh-bridge -c examples/ghostify.ssh.json
```
