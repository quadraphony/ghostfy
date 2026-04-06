#!/usr/bin/env bash
set -euo pipefail

BUILD_DIR="${BUILD_DIR:-build}"
mkdir -p "$BUILD_DIR"

VERSION="$(git describe --always --dirty="-dirty" || echo "dev")"
echo "{\"version\":\"$VERSION\",\"timestamp\":\"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"}" > "$BUILD_DIR/version.json"

GO_BIN="${GO_BIN:-/snap/go/current/bin/go}"
GOOS=linux GOARCH=amd64 "$GO_BIN" build -o "$BUILD_DIR/ghostify-linux" ./cmd/ghostify
GOOS=windows GOARCH=amd64 "$GO_BIN" build -o "$BUILD_DIR/ghostify-windows.exe" ./cmd/ghostify
