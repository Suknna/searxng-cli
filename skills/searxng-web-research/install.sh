#!/usr/bin/env bash
set -euo pipefail

SKILL_NAME="searxng-web-research"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TARGET_BASE="${TARGET_SKILLS_DIR:-$HOME/.config/opencode/skills}"
TARGET_DIR="$TARGET_BASE/$SKILL_NAME"

mkdir -p "$TARGET_BASE"
rm -rf "$TARGET_DIR"
cp -R "$SCRIPT_DIR" "$TARGET_DIR"

printf "Installed skill to %s\n" "$TARGET_DIR"
printf "Next: ensure 'searxng-cli' is available on PATH.\n"
