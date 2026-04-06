# UDP2RAW Bridge

Run UDP2RAW in client or server mode to encapsulate traffic into a stealth transport.

## Example

```json
{
  "profile": "udp2raw-bridge",
  "log_level": "info",
  "args": [
    "-c",
    "-l",
    "0.0.0.0:4000",
    "-r",
    "51.68.155.153:4000",
    "--raw-mode",
    "faketcp"
  ]
}
```

Run it via:

```
ghostify udp2raw-bridge -c examples/ghostify.udp2raw.json
```
