# searxng-cli

> [English](README.md) | [中文](README_CN.md)

A CLI tool optimized for LLMs and AI agents. Search the web and extract page content in a compact, context-efficient format.

## Overview

`searxng-cli` provides a two-step workflow for web research:

1. **`search`** - Discover URLs via SearXNG, output a compact Markdown table
2. **`read`** - Extract clean page content (markdown or text) without browser automation

**Core rule:** Search first, then read. If you already know the URL, run `read` directly.

## Prerequisites

You need access to a SearXNG instance:

- **Option A:** Self-deploy SearXNG ([official docs](https://docs.searxng.org/admin/installation.html))
- **Option B:** Use a public/community instance (e.g., `https://search.sapti.me`)

## Installation For Humans

Copy and paste this prompt to your LLM agent (Claude Code, OpenCode, Codex, Cursor, etc.):

```text
Install and configure searxng-cli by following the instructions here:
https://raw.githubusercontent.com/Suknna/searxng-cli/refs/heads/main/install.md
```

Or read `install.md` and do it manually, but seriously, let an agent do it. Humans fat-finger configs.

If you need custom behavior (non-default paths, enterprise policy, custom auth flow), switch to manual installation and adapt commands from `install.md`.

## Quick Start

### 1. Install

```bash
# Download from releases
curl -L -o searxng-cli https://github.com/your-org/searxng-cli/releases/latest/download/searxng-cli_linux_amd64
chmod +x searxng-cli
sudo mv searxng-cli /usr/local/bin/
```

Or build from source:

```bash
go build -o searxng-cli .
sudo mv searxng-cli /usr/local/bin/
```

### 2. Initialize Configuration

Create a permanent config file:

```bash
searxng-cli config init
```

This generates `~/.config/searxng-cli/config.yml` with default settings.

**Configuration file location:**
- Linux/macOS: `~/.config/searxng-cli/config.yml`
- Windows: `%APPDATA%\searxng-cli\config.yml`

### 3. Set Your SearXNG Server

Edit `~/.config/searxng-cli/config.yml`:

```yaml
apiVersion: searxng-cli/v1
kind: Config
current-context: default
contexts:
  default:
    server: "https://your-searxng-instance.com/"
    timeout: "10s"
    limit: 10
```

Or use a one-liner to view/update:

```bash
searxng-cli config view          # View effective config
searxng-cli config use-context <name>  # Switch context
```

### 4. Start Searching

```bash
# Search for information
searxng-cli search "golang context cancellation best practices"

# Extract content from a result
searxng-cli read "https://go.dev/blog/context"
```

## Commands

### `search <query>`

Search SearXNG and output a Markdown table.

```bash
searxng-cli search "machine learning" --limit 5
```

**Output format:**
```markdown
# | title | url | content | template
| 1 | Title Here | https://example.com | Summary... | Title=Title Here URL=https://example.com Content=Summary...
```

**Key flags:**
- `--limit <n>`: Limit results (default: 10)
- `--template <string>`: Custom output template
- `--server <url>`: Override SearXNG server
- `--timeout <duration>`: Request timeout

### `read <url>`

Extract clean page content without browser automation.

```bash
searxng-cli read "https://go.dev/blog/context" --format markdown
searxng-cli read "https://go.dev/blog/context" --format text
```

**Key flags:**
- `--format <markdown|text>`: Output format (default: markdown)
- `--timeout <duration>`: Request timeout (default: 10s)
- `--respect-robots <true|false>`: Check robots.txt (default: true)
- `--max-bytes <n>`: Max response size (default: 2MB)
- `--retry <n>`: Retry count (default: 1)

### Configuration Management

```bash
searxng-cli config init                 # Create default config
searxng-cli config view                 # Show effective config
searxng-cli config use-context <name>   # Switch active context
```

## Configuration Precedence

Settings are resolved in this order (highest first):

1. **Command-line flags** - `--server`, `--timeout`, etc.
2. **Environment variables** - `SEARXNG_CLI_*`
3. **Configuration file** - `~/.config/searxng-cli/config.yml`
4. **Built-in defaults**

## Why searxng-cli?

### Compared to SearXNG MCP / Browser Automation

| Aspect | searxng-cli | SearXNG MCP / Browser Tools |
|--------|-------------|----------------------------|
| **Context Length** | Compact markdown tables and cleaned text | Full HTML, often with scripts and noise |
| **Output Format** | Only what LLMs need: titles, URLs, clean content | Raw or minimally processed HTML |
| **Resource Usage** | No browser, lightweight HTTP requests | Requires browser engine (Chrome, Playwright) |
| **Speed** | Fast HTTP calls | Slower (DOM rendering, JS execution) |
| **Use Case** | Research, summarization, citation | Interactive navigation, form submission, visual testing |

**Key advantage:** `searxng-cli` is purpose-built for LLM workflows. It strips away visual layout, scripts, and styling—delivering only the semantic content your model needs to analyze and respond.

### Best For

- Research tasks requiring clean text extraction
- Building citation lists
- Multi-source comparison
- Content summarization
- Automated documentation updates

**Not for:** Login flows, form submission, JavaScript-heavy interactions (use browser automation instead).

## Skill for AI Agents

This repository includes an installable skill for OpenCode, Claude, and other AI agents:

```bash
# Install the skill
bash skills/searxng-web-research/install.sh
```

The skill enables agents to:
- Search for sources and validate URLs
- Extract page content in a format optimized for LLM context windows
- Perform multi-source research with clean, deduplicated content

## Authentication

If your SearXNG instance requires authentication:

```bash
# API Key mode
searxng-cli --auth-mode api_key \
  --auth-header "X-API-Key" \
  --auth-api-key "$(echo -n 'your-key' | base64)" \
  search "query"

# Basic auth mode
searxng-cli --auth-mode basic \
  --auth-username "$(echo -n 'user' | base64)" \
  --auth-password "$(echo -n 'pass' | base64)" \
  search "query"
```

Or configure permanently in `config.yml`:

```yaml
contexts:
  default:
    server: "https://your-instance.com/"
    auth:
      mode: "api_key"
      api_key_header: "X-API-Key"
      api_key: "<base64-encoded-key>"
```

## Error Handling

Errors are written to `stderr` as structured `key=value` pairs:

```
code=NETWORK_TIMEOUT message="network request timed out" retryable=true hint="check network path or increase timeout"
```

Common error codes:
- `NETWORK_TIMEOUT`: Increase `--timeout` or retry
- `HTTP_NON_2XX`: Check server health and URL
- `ROBOTS_DISALLOWED`: Site blocks crawling (use `--respect-robots=false` with caution)
- `EXTRACT_EMPTY`: Page has no extractable content

## Development

```bash
# Run tests
go test ./...

# Run a specific test
go test ./cmd -run '^TestReadHelp$' -v

# Format code
gofmt -w $(rg --files -g '*.go')
go vet ./...

# Build locally
go build -o searxng-cli .
```

See [AGENTS.md](./AGENTS.md) for coding conventions and agent workflow guidelines.

## License

MIT
