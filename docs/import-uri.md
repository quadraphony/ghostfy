# Import URI

Ghostify can convert a supported share URI into Ghostify JSON.

## Supported Today

- `vless://`
- `vmess://`
- `trojan://`

## Example

```bash
ghostify import-uri 'vless://43488128-319e-f480-64ea-0acdc712e2a8@51.68.155.153:443?encryption=none&flow=xtls-rprx-vision&fp=chrome&pbk=_xzl59bcUYD9QUSKFyboqCC_9eUnaXUUKWA19oQWFHU&security=reality&sid=8d17b7ebab7b99d7&sni=www.allegro.pl&type=tcp#PL-vless'
```

The output is Ghostify-owned JSON that can be saved and used with:

```bash
ghostify run -c ghostify.json
```
