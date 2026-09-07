#!/usr/bin/env bash
set -euo pipefail

BIN_NAME="madv"
DEST_DIR="${BINDIR:-/usr/bin}"
TARGET="$DEST_DIR/$BIN_NAME"

if [[ ! -f "$TARGET" ]]; then
    echo "Notice: $TARGET does not exist. Nothing to uninstall."
    exit 0
fi

echo "==> Removing $TARGET..."
SUDO=""
if [[ "$(id -u)" -ne 0 && ! -w "$DEST_DIR" ]]; then
    SUDO="sudo"
fi

$SUDO rm -f "$TARGET"

if [[ ! -f "$TARGET" ]]; then
    echo "==> Successfully uninstalled $BIN_NAME from $DEST_DIR"
else
    echo "Error: Failed to remove $TARGET" >&2
    exit 1
fi
