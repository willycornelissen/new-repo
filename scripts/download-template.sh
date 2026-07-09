#!/usr/bin/env bash
set -euo pipefail

REPO_URL="https://github.com/willycornelissen/ai-template"
TEMPLATE_DIR="$(cd "$(dirname "$0")/.." && pwd)/internal/embed/template"

echo "Downloading template from $REPO_URL ..."

TMP_DIR=$(mktemp -d)
trap 'rm -rf "$TMP_DIR"' EXIT

git clone --depth 1 "$REPO_URL" "$TMP_DIR"

rm -rf "$TEMPLATE_DIR"
mkdir -p "$TEMPLATE_DIR"

shopt -s dotglob
for item in "$TMP_DIR"/*; do
  base=$(basename "$item")
  if [ "$base" = ".git" ]; then
    continue
  fi
  # go:embed does not support directories starting with '.', so rename .opencode -> opencode
  if [ "$base" = ".opencode" ]; then
    cp -r "$item" "$TEMPLATE_DIR/opencode"
  else
    cp -r "$item" "$TEMPLATE_DIR/"
  fi
done
shopt -u dotglob

echo "Template updated at $TEMPLATE_DIR"
echo "Size: $(du -sh "$TEMPLATE_DIR" | cut -f1)"
