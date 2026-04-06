# Packaging

Run `scripts/build.sh` to produce Linux and Windows binaries under `build/` plus `build/version.json`.

```
./scripts/build.sh
```

- For releases, bundle the generated binaries with `build/version.json`.
- Update `build/version.json` before tagging if you need a semantic version.
- Installation is as simple as copying the binary to a location like `/usr/local/bin` or using the Windows `.exe`.
