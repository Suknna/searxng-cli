# searxng-web-research skill

This directory contains an installable skill for agents that need web search and page reading through `searxng-cli`.

## What this skill provides

- Web discovery with `searxng-cli search`
- Page extraction with `searxng-cli read`
- Guidance for non-interactive research flows

## Install (OpenCode-compatible)

From the repository root:

```bash
bash skills/searxng-web-research/install.sh
```

Optional custom target:

```bash
TARGET_SKILLS_DIR="$HOME/.config/opencode/skills" bash skills/searxng-web-research/install.sh
```

## Requirements

- `searxng-cli` must be available in PATH, or the agent must run commands via `go run .`
- Network access to the configured SearXNG server and target websites

## Verify

```bash
searxng-cli search --help
searxng-cli read --help
```

Then inspect `SKILL.md` in the installed location.
