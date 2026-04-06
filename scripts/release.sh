#!/usr/bin/env bash
set -euo pipefail

BUILD_DIR="${BUILD_DIR:-build}"
RELEASE_DIR="${RELEASE_DIR:-release}"
mkdir -p "$BUILD_DIR" "$RELEASE_DIR"

./scripts/build.sh

VERSION_FILE="$BUILD_DIR/version.json"
VERSION="$(python3 -c 'import json,sys; print(json.load(open(sys.argv[1]))["version"])' "$VERSION_FILE")"

tar -C "$BUILD_DIR" -czf "$RELEASE_DIR/ghostify-${VERSION}.tar.gz" ghostify-linux ghostify-windows.exe version.json
zip -j "$RELEASE_DIR/ghostify-${VERSION}.zip" "$BUILD_DIR"/ghostify-linux "$BUILD_DIR"/ghostify-windows.exe "$VERSION_FILE"
