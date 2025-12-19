#!/usr/bin/env bash
set -euo pipefail

BUILD_DIR="./dist"
SHASUM_TXT="$BUILD_DIR/shasum.txt"

mkdir -p "$BUILD_DIR"

builds=(
  "darwin amd64 lista-darwin-amd64"
  "darwin arm64 lista-darwin-arm64"
  "linux  amd64 lista-linux-amd64"
  "linux  arm64 lista-linux-arm64"
)

echo "Building binaries..."

for build in "${builds[@]}"; do
  read -r GOOS GOARCH OUTPUT <<< "$build"

  echo "→ $OUTPUT"
  GOOS=$GOOS GOARCH=$GOARCH \
    go build -o "$BUILD_DIR/$OUTPUT"
done

#File is created if missing, File is emptied if it exists
: > "$SHASUM_TXT"

echo "Generating checksums..."

for build_file in "$BUILD_DIR"/lista-*; do
  shasum -a 256 -b "$build_file" >> "$SHASUM_TXT"
done

echo "Done."
