# searxng-cli

> [English](README.md) | [中文](README_CN.md)

为 LLM 和 AI 智能体优化的 CLI 工具。以紧凑、节省上下文的格式搜索网页并提取页面内容。

## 概述

`searxng-cli` 提供两步式网页研究工作流：

1. **`search`** - 通过 SearXNG 发现 URL，输出紧凑的 Markdown 表格
2. **`read`** - 提取干净的页面内容（markdown 或文本），无需浏览器自动化

**核心规则：** 先搜索，再阅读。如果你已经知道 URL，直接运行 `read`。

## 前置条件

你需要访问 SearXNG 实例：

- **选项 A：** 自行部署 SearXNG ([官方文档](https://docs.searxng.org/admin/installation.html))
- **选项 B：** 使用公开/社区实例（例如 `https://search.sapti.me`）

## 快速开始

### 1. 安装

```bash
# 从 releases 下载
curl -L -o searxng-cli https://github.com/your-org/searxng-cli/releases/latest/download/searxng-cli_linux_amd64
chmod +x searxng-cli
sudo mv searxng-cli /usr/local/bin/
```

或者从源码构建：

```bash
go build -o searxng-cli .
sudo mv searxng-cli /usr/local/bin/
```

### 2. 初始化配置

创建永久配置文件：

```bash
searxng-cli config init
```

这会生成包含默认设置的 `~/.config/searxng-cli/config.yml`。

**配置文件位置：**
- Linux/macOS: `~/.config/searxng-cli/config.yml`
- Windows: `%APPDATA%\searxng-cli\config.yml`

### 3. 设置 SearXNG 服务器

编辑 `~/.config/searxng-cli/config.yml`：

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

或者使用单行命令查看/更新：

```bash
searxng-cli config view          # 查看有效配置
searxng-cli config use-context <name>  # 切换上下文
```

### 4. 开始搜索

```bash
# 搜索信息
searxng-cli search "golang context cancellation best practices"

# 从结果中提取内容
searxng-cli read "https://go.dev/blog/context"
```

## 命令

### `search <query>`

搜索 SearXNG 并输出 Markdown 表格。

```bash
searxng-cli search "machine learning" --limit 5
```

**输出格式：**
```markdown
# | title | url | content | template
| 1 | 标题 | https://example.com | 摘要... | Title=标题 URL=https://example.com Content=摘要...
```

**关键参数：**
- `--limit <n>`: 限制结果数量（默认：10）
- `--template <string>`: 自定义输出模板
- `--server <url>`: 覆盖 SearXNG 服务器
- `--timeout <duration>`: 请求超时时间

### `read <url>`

无需浏览器自动化提取干净的页面内容。

```bash
searxng-cli read "https://go.dev/blog/context" --format markdown
searxng-cli read "https://go.dev/blog/context" --format text
```

**关键参数：**
- `--format <markdown|text>`: 输出格式（默认：markdown）
- `--timeout <duration>`: 请求超时时间（默认：10s）
- `--respect-robots <true|false>`: 检查 robots.txt（默认：true）
- `--max-bytes <n>`: 最大响应大小（默认：2MB）
- `--retry <n>`: 重试次数（默认：1）

### 配置管理

```bash
searxng-cli config init                 # 创建默认配置
searxng-cli config view                 # 显示有效配置
searxng-cli config use-context <name>   # 切换活动上下文
```

## 配置优先级

设置按以下顺序解析（优先级从高到低）：

1. **命令行参数** - `--server`, `--timeout` 等
2. **环境变量** - `SEARXNG_CLI_*`
3. **配置文件** - `~/.config/searxng-cli/config.yml`
4. **内置默认值**

## 为什么选择 searxng-cli？

### 与 SearXNG MCP / 浏览器自动化对比

| 方面 | searxng-cli | SearXNG MCP / 浏览器工具 |
|------|-------------|------------------------|
| **上下文长度** | 紧凑的 markdown 表格和清理后的文本 | 完整 HTML，通常包含脚本和噪声 |
| **输出格式** | 仅 LLM 需要的内容：标题、URL、干净内容 | 原始或最小处理的 HTML |
| **资源使用** | 无需浏览器，轻量级 HTTP 请求 | 需要浏览器引擎（Chrome、Playwright）|
| **速度** | 快速 HTTP 调用 | 较慢（DOM 渲染、JS 执行）|
| **使用场景** | 研究、摘要、引用 | 交互式导航、表单提交、视觉测试 |

**核心优势：** `searxng-cli` 专为 LLM 工作流打造。它剥离视觉布局、脚本和样式，只提供模型分析所需的语义内容。

### 最适合

- 需要干净文本提取的研究任务
- 构建引用列表
- 多来源对比
- 内容摘要
- 自动化文档更新

**不适合：** 登录流程、表单提交、重度 JavaScript 交互（请改用浏览器自动化）。

## 用于 AI 智能体的 Skill

本仓库包含适用于 OpenCode、Claude 和其他 AI 智能体的可安装 skill：

```bash
# 安装 skill
bash skills/searxng-web-research/install.sh
```

该 skill 使智能体能够：
- 搜索来源并验证 URL
- 以针对 LLM 上下文窗口优化的格式提取页面内容
- 使用干净、去重后的内容执行多来源研究

## 认证

如果你的 SearXNG 实例需要认证：

```bash
# API Key 模式
searxng-cli --auth-mode api_key \
  --auth-header "X-API-Key" \
  --auth-api-key "$(echo -n 'your-key' | base64)" \
  search "query"

# Basic 认证模式
searxng-cli --auth-mode basic \
  --auth-username "$(echo -n 'user' | base64)" \
  --auth-password "$(echo -n 'pass' | base64)" \
  search "query"
```

或者在 `config.yml` 中永久配置：

```yaml
contexts:
  default:
    server: "https://your-instance.com/"
    auth:
      mode: "api_key"
      api_key_header: "X-API-Key"
      api_key: "<base64-encoded-key>"
```

## 错误处理

错误以结构化 `key=value` 对形式写入 `stderr`：

```
code=NETWORK_TIMEOUT message="network request timed out" retryable=true hint="check network path or increase timeout"
```

常见错误代码：
- `NETWORK_TIMEOUT`: 增加 `--timeout` 或重试
- `HTTP_NON_2XX`: 检查服务器健康和 URL
- `ROBOTS_DISALLOWED`: 站点阻止爬取（谨慎使用 `--respect-robots=false`）
- `EXTRACT_EMPTY`: 页面没有可提取内容

## 开发

```bash
# 运行测试
go test ./...

# 运行特定测试
go test ./cmd -run '^TestReadHelp$' -v

# 格式化代码
gofmt -w $(rg --files -g '*.go')
go vet ./...

# 本地构建
go build -o searxng-cli .
```

查看 [AGENTS.md](./AGENTS.md) 了解编码规范和智能体工作流指南。

## 许可证

MIT
