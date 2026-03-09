---
name: searxng-web-research
description: Use when an agent needs lightweight web research by searching with searxng-cli and reading webpage content without browser automation.
allowed-tools: Bash(searxng-cli:*), Bash(go:*)
---

# SearXNG Web Research

## Overview

This skill gives agents a two-step workflow:
1) use `search` to discover candidate URLs, 2) use `read` to extract clean page content.

Core rule: search first, then read. If the URL is already known, run `read` directly.

## When to Use

- Need quick source discovery from the open web
- Need clean article/page text for summarization or comparison
- Need non-interactive retrieval (no login, no clicking UI flows)

Do not use this skill when:
- The task requires login, form submission, or interactive navigation
- The content only exists after heavy client-side JavaScript rendering

## Quick Start

If `searxng-cli` is on PATH:

```bash
searxng-cli search "golang context cancellation best practices"
searxng-cli read "https://go.dev/blog/context"
```

If running from this repository without installing the binary:

```bash
go run . search "golang context cancellation best practices"
go run . read "https://go.dev/blog/context"
```

## Core Workflow

### 1) Search phase (discover candidates)

```bash
searxng-cli search "<query>"
```

Common flags:
- `--limit <n>`: local output row limit (default 10)
- `--template "..."`: custom per-row template

Output is a Markdown table with fixed columns:
- `# | title | url | content | template`

### 2) Read phase (extract target page)

```bash
searxng-cli read "<url>"
```

Common flags:
- `--format markdown|text` (default `markdown`)
- `--timeout 10s`
- `--respect-robots=true|false`
- `--max-bytes 2097152`
- `--retry 1`

## Recommended Patterns

- Known URL: skip search and run `read` directly
- Multi-source validation:
  1) collect 3-5 candidate URLs with `search`
  2) extract each page via `read --format text`
  3) deduplicate claims and cross-check sources
- Prefer high-quality sources: official docs, standards, and reputable publishers

## Auth and Config

For authenticated SearXNG instances, pass global auth flags (values are base64-encoded):

```bash
searxng-cli --auth-mode api_key --auth-header "X-API-Key" --auth-api-key "<base64>" search "<query>"
```

Config precedence is stable: `flags > env > config > defaults`.

## Failure Handling

Errors are written to `stderr` as structured `key=value` lines.

Typical cases:
- `NETWORK_TIMEOUT`: increase `--timeout` or retry later
- `HTTP_NON_2XX`: validate upstream health, URL, and auth config
- `ROBOTS_DISALLOWED`: follow site policy; disable robots checks only when explicitly allowed
- `EXTRACT_EMPTY`: try another URL or compare with `--format text`

## Verification Checklist

- `searxng-cli search --help` shows search usage
- `searxng-cli read --help` shows read usage
- `searxng-cli search "test"` prints a Markdown result table
- `searxng-cli read "https://example.com"` returns readable content

## Security Notes

- Do not include secrets in CLI arguments unless required.
- Keep robots checks enabled by default.
- This skill is intentionally non-interactive; use browser automation skills for UI workflows.
