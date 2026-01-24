#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(dirname "$SCRIPT_DIR")"
HOOKS_DIR="$REPO_ROOT/.git/hooks"

echo "Setting up git hooks..."

# Create hooks directory if it doesn't exist
mkdir -p "$HOOKS_DIR"

# Create symlink for pre-push
ln -sf "$SCRIPT_DIR/pre-push" "$HOOKS_DIR/pre-push"
chmod +x "$SCRIPT_DIR/pre-push"

echo "Done! Pre-push hook installed."
