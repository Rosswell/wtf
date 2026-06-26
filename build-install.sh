#!/usr/bin/env bash
# Build this wtf fork and install it on PATH via a symlink in ~/.local/bin.
#
# ~/.local/bin is owned by us (not brew), so the symlink survives brew
# operations. Re-run after editing any module to rebuild + relink.
set -euo pipefail

REPO_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BIN="$REPO_DIR/bin/wtfutil"
TARGET="$HOME/.local/bin/wtfutil"

echo "Building wtf fork in $REPO_DIR ..."
cd "$REPO_DIR"
go build -o "$BIN" .

echo "Linking $TARGET -> $BIN"
ln -sf "$BIN" "$TARGET"

echo "Done. $("$TARGET" --version 2>/dev/null | head -1)"
echo "Custom modules: selectablecmd"
