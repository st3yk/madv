#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

BIN_NAME="madv"
DEST_DIR="${BINDIR:-/usr/bin}"

echo "==> Building $BIN_NAME with Bazel..."
bazel build //:madv

BUILT_BIN="bazel-bin/cmd/madv/madv_/$BIN_NAME"
if [[ ! -f "$BUILT_BIN" ]]; then
    echo "Error: Built binary not found at $BUILT_BIN" >&2
    exit 1
fi

echo "==> Installing $BIN_NAME to $DEST_DIR..."
SUDO=""
if [[ "$(id -u)" -ne 0 && ! -w "$DEST_DIR" ]]; then
    SUDO="sudo"
fi

$SUDO install -m 755 "$BUILT_BIN" "$DEST_DIR/$BIN_NAME"

echo "==> Successfully installed $BIN_NAME to $DEST_DIR/$BIN_NAME"
if [[ -x "$DEST_DIR/$BIN_NAME" ]]; then
    echo "==> Installed version: $("$DEST_DIR/$BIN_NAME" --version)"
fi
