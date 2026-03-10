# searxng-cli 安装与 Agent 接入指南

本文档面向 AI Agent 执行安装任务，目标是：
- 安装 `searxng-cli` 二进制
- 完成用户本地永久配置
- 验证 `search/read` 可用
- 将本项目 skill 安装到目标 Agent 的默认 skill 目录

## Agent 安装 TODO List（必须逐项勾选）

- [ ] 1. 识别用户操作系统与 CPU 架构
- [ ] 2. 从 GitHub Release 下载匹配的二进制文件
- [ ] 3. 安装到 `/usr/local/bin/searxng-cli`
- [ ] 4. 刷新 shell 环境并验证 `searxng-cli --help`
- [ ] 5. 收集用户 SearXNG 配置信息并写入 `~/.config/searxng-cli/config.yml`
- [ ] 6. 执行可用性验证：`search "hello world"` + `read <第一个URL>`
- [ ] 7. 安装 skill 到用户 Agent 默认 skill 目录（OpenCode/Claude Code/Codex）
- [ ] 8. 通知用户重启 Agent 客户端以重新加载 skill

## 1) 项目简介

`searxng-cli` 是一个面向大模型工作流的搜索与网页读取 CLI：
- `search`：通过 SearXNG 获取搜索结果，并输出紧凑 Markdown 表格
- `read`：读取网页正文并做净化，仅保留模型需要的信息

相对 SearXNG MCP / 浏览器自动化工具，本项目重点优势：
- 缩短上下文：输出结构化且紧凑，减少噪声 token
- 结果可直接喂给模型：标题、URL、摘要、正文均已格式化
- 无浏览器依赖：更轻量，执行路径更短

## 2) 安装前提

用户必须具备以下条件：
- 可访问一个 SearXNG 实例（自建或社区公开实例）
- 终端具备 `curl`、`chmod`、`mv` 等基础命令
- 具备写入 `/usr/local/bin` 权限（通常需要 sudo）

Windows 说明：
- 当前项目**未在 Windows 上做过完整测试**
- 若用户是 Windows，请让用户自行拉取仓库后本地编译安装
- 不要承诺 Windows 预编译包可直接使用

## 3) 安装步骤（二进制方式）

### 步骤 1：检测操作系统与架构

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

### 步骤 2：选择 release 资产并下载

仓库地址：`https://github.com/Suknna/searxng-cli`

当前 release workflow 产物：
- `searxng-cli_linux_amd64`
- `searxng-cli_linux_arm64`
- `searxng-cli_darwin_arm64`

下载命令（自动拼接）：

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

### 步骤 3：移动到 `/usr/local/bin`

```bash
sudo mv searxng-cli /usr/local/bin/searxng-cli
```

### 步骤 4：刷新环境变量并验证

```bash
hash -r
searxng-cli --help
```

如果命令找不到：
- 确认 `/usr/local/bin` 在 `PATH` 中
- 重新打开一个终端会话后再试

## 4) 生成并写入本地永久配置

### 步骤 5.1：询问用户配置（最少实例 URL）

最少必须问到：
- `server`（SearXNG 实例 URL，必填）

可选项（可留空）：
- `auth_mode`：`none` / `api_key` / `basic`
- `auth_header`（api_key 模式）
- `auth_api_key`（base64）
- `auth_username`（base64，basic 模式）
- `auth_password`（base64，basic 模式）
- `timeout`（默认 `10s`）
- `limit`（默认 `10`）

### 步骤 5.2：初始化默认配置

```bash
searxng-cli config init
```

默认配置路径：
- macOS/Linux：`~/.config/searxng-cli/config.yml`

### 步骤 5.3：根据用户回答修改 YAML

建议直接覆盖写入（不要半改导致结构错乱）：

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

## 5) 配置后验证（必须执行）

### 步骤 6：执行搜索与网页读取

1. 运行搜索：

```bash
searxng-cli search "hello world"
```

2. 从结果中取第一个 URL，再读取正文：

```bash
# 这里的 FIRST_URL 由 agent 从 search 输出中提取
searxng-cli read "${FIRST_URL}"
```

判定标准：
- `search` 成功返回表格
- `read` 成功返回正文（markdown 或 text）

## 6) 安装 skill 到 Agent 默认目录

注意：`skills/searxng-web-research` 是**远端仓库内路径**，不是用户本地默认已存在的路径。

远端仓库：`https://github.com/Suknna/searxng-cli`

先从远端拉取 skill 目录到本地临时目录：

```bash
TMP_DIR="$(mktemp -d)"
git clone --depth 1 https://github.com/Suknna/searxng-cli.git "$TMP_DIR/searxng-cli"
REMOTE_SKILL_DIR="$TMP_DIR/searxng-cli/skills/searxng-web-research"

test -f "$REMOTE_SKILL_DIR/SKILL.md" || { echo "SKILL.md not found in remote repo"; exit 1; }
```

### OpenCode

默认全局目录：`~/.config/opencode/skills`

```bash
mkdir -p "$HOME/.config/opencode/skills"
cp -R "$REMOTE_SKILL_DIR" "$HOME/.config/opencode/skills/"
```

参考：`https://opencode.ai/docs/skills/`

### Claude Code

默认全局目录：`~/.claude/skills`

```bash
mkdir -p "$HOME/.claude/skills"
cp -R "$REMOTE_SKILL_DIR" "$HOME/.claude/skills/"
```

参考：`https://code.claude.com/docs/en/skills`

### Codex

默认全局目录：`~/.agents/skills`

```bash
mkdir -p "$HOME/.agents/skills"
cp -R "$REMOTE_SKILL_DIR" "$HOME/.agents/skills/"
```

参考：`https://developers.openai.com/codex/skills/`

### 环境清理

安装完成后可清理临时目录：

```bash
rm -rf "$TMP_DIR"
```

## 7) 完成提示（必须告知用户）

步骤 7：安装 skill 后，明确提示用户：

- 请重启 Agent 客户端（或新开会话）
- skill 列表会在新会话中重新加载
- 若未显示 skill，优先检查目标路径下是否存在 `SKILL.md`

## 8) 故障排查

- 下载失败：检查网络与 GitHub 访问权限
- 命令不可用：检查 `/usr/local/bin` 是否在 `PATH`
- `search` 失败：优先检查 `server` URL 可达性
- `read` 失败：检查目标 URL 可访问性、robots 策略、超时配置
- skill 不生效：确认目录名与 `SKILL.md` frontmatter 的 `name` 一致
