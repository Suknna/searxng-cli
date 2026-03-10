# searxng-cli Installation and Agent Integration Guide

This document is intended for AI Agents performing installation tasks. The objectives are:
- Install the `searxng-cli` binary
- Complete user's local permanent configuration
- Verify `search/read` functionality
- Install this project's skill to the target Agent's default skill directory

## Agent Installation TODO List (Must check each item)

- [ ] 1. Identify user's operating system and CPU architecture
- [ ] 2. Download matching binary from GitHub Release
- [ ] 3. Install to `/usr/local/bin/searxng-cli`
- [ ] 4. Refresh shell environment and verify `searxng-cli --help`
- [ ] 5. Collect user's SearXNG configuration and write to `~/.config/searxng-cli/config.yml`
- [ ] 6. Execute usability verification: `search "hello world"` + `read <first URL>`
- [ ] 7. Install skill to user's Agent default skill directory (OpenCode/Claude Code/Codex)
- [ ] 8. Notify user to restart Agent client to reload skill

## 1) Project Introduction

`searxng-cli` is a search and web page reading CLI designed for LLM workflows:
- `search`: Get search results via SearXNG and output compact Markdown tables
- `read`: Read web page content and sanitize it, keeping only the information the model needs

Key advantages over SearXNG MCP / browser automation tools:
- Reduced context: Structured and compact output, fewer noise tokens
- Results can be directly fed to the model: title, URL, summary, and content are all formatted
- No browser dependency: Lighter weight, shorter execution path

## 2) Prerequisites

Users must have the following:
- Access to a SearXNG instance (self-hosted or community public instance)
- Terminal with basic commands like `curl`, `chmod`, `mv`
- Permission to write to `/usr/local/bin` (usually requires sudo)

Windows Note:
- This project has **not been fully tested on Windows**
- If the user is on Windows, ask them to clone the repository and build locally
- Do not promise that Windows prebuilt binaries will work directly

## 3) Installation Steps (Binary Method)

### Step 1: Detect Operating System and Architecture

```bash
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case "$ARCH" in
  x86_64) ARCH="amd64" ;;
  arm64|aarch64) ARCH="arm64" ;;
  *) echo "Unsupported arch: $ARCH"; exit 1 ;;
esac

echo "Detected OS=$OS ARCH=$ARCH"
```

### Step 2: Select release asset and download

Repository URL: `https://github.com/Suknna/searxng-cli`

Current release workflow artifacts:
- `searxng-cli_linux_amd64`
- `searxng-cli_linux_arm64`
- `searxng-cli_darwin_arm64`

Download command (auto-assembled):

```bash
case "$OS-$ARCH" in
  linux-amd64|linux-arm64|darwin-arm64)
    ASSET="searxng-cli_${OS}_${ARCH}"
    ;;
  *)
    echo "No prebuilt binary for $OS-$ARCH. Please build from source."; exit 1
    ;;
esac

curl -fL "https://github.com/Suknna/searxng-cli/releases/latest/download/${ASSET}" -o searxng-cli
chmod +x searxng-cli
```

### Step 3: Move to `/usr/local/bin`

```bash
sudo mv searxng-cli /usr/local/bin/searxng-cli
```

### Step 4: Refresh environment variables and verify

```bash
hash -r
searxng-cli --help
```

If command not found:
- Confirm `/usr/local/bin` is in `PATH`
- Try reopening a terminal session

## 4) Generate and Write Local Permanent Configuration

### Step 5.1: Ask user for configuration (minimum: instance URL)

Minimum required:
- `server` (SearXNG instance URL, required)

Optional fields (can be left empty):
- `auth_mode`: `none` / `api_key` / `basic`
- `auth_header` (api_key mode)
- `auth_api_key` (base64)
- `auth_username` (base64, basic mode)
- `auth_password` (base64, basic mode)
- `timeout` (default `10s`)
- `limit` (default `10`)

### Step 5.2: Initialize default configuration

```bash
searxng-cli config init
```

Default configuration path:
- macOS/Linux: `~/.config/searxng-cli/config.yml`

### Step 5.3: Modify YAML based on user's answers

Recommend overwriting directly (avoid partial edits that may corrupt structure):

```yaml
apiVersion: searxng-cli/v1
kind: Config
current-context: default
contexts:
  default:
    server: "https://your-searxng-instance.example/"
    timeout: "10s"
    limit: 10
    template: "Title={{.Title}} URL={{.URL}} Content={{.Content}}"
    auth:
      mode: "none"
      api_key_header: "X-API-Key"
      api_key: ""
      username: ""
      password: ""
```

## 5) Post-Configuration Verification (Must Execute)

### Step 6: Execute search and web page reading

1. Run search:

```bash
searxng-cli search "hello world"
```

2. Extract the first URL from results, then read content:

```bash
# FIRST_URL should be extracted by agent from search output
searxng-cli read "${FIRST_URL}"
```

Success criteria:
- `search` successfully returns a table
- `read` successfully returns content (markdown or text)

## 6) Install skill to Agent default directory

Note: `skills/searxng-web-research` is a **remote repository path**, not a local default path.

Remote repository: `https://github.com/Suknna/searxng-cli`

First, clone the skill directory from remote to a local temporary directory:

```bash
TMP_DIR="$(mktemp -d)"
git clone --depth 1 https://github.com/Suknna/searxng-cli.git "$TMP_DIR/searxng-cli"
REMOTE_SKILL_DIR="$TMP_DIR/searxng-cli/skills/searxng-web-research"

test -f "$REMOTE_SKILL_DIR/SKILL.md" || { echo "SKILL.md not found in remote repo"; exit 1; }
```

### OpenCode

Default global directory: `~/.config/opencode/skills`

```bash
mkdir -p "$HOME/.config/opencode/skills"
cp -R "$REMOTE_SKILL_DIR" "$HOME/.config/opencode/skills/"
```

Reference: `https://opencode.ai/docs/skills/`

### Claude Code

Default global directory: `~/.claude/skills`

```bash
mkdir -p "$HOME/.claude/skills"
cp -R "$REMOTE_SKILL_DIR" "$HOME/.claude/skills/"
```

Reference: `https://code.claude.com/docs/en/skills`

### Codex

Default global directory: `~/.agents/skills`

```bash
mkdir -p "$HOME/.agents/skills"
cp -R "$REMOTE_SKILL_DIR" "$HOME/.agents/skills/"
```

Reference: `https://developers.openai.com/codex/skills/`

### Environment Cleanup

After installation, clean up temporary directory:

```bash
rm -rf "$TMP_DIR"
```

## 7) Completion Notice (Must notify user)

Step 7: After installing the skill, clearly notify the user:

- Please restart the Agent client (or start a new session)
- The skill list will be reloaded in the new session
- If the skill does not appear, first check if `SKILL.md` exists in the target path

## 8) Troubleshooting

- Download failed: Check network and GitHub access permissions
- Command not available: Check if `/usr/local/bin` is in `PATH`
- `search` failed: First check `server` URL accessibility
- `read` failed: Check target URL accessibility, robots policy, timeout configuration
- Skill not working: Confirm directory name matches `SKILL.md` frontmatter `name` field