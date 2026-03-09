# searxng-cli MVP 规格说明（面向 LLM）

## 目的
一个用于 LLM/智能体的 Go（Cobra）CLI，包含两类能力：
- `search`：通过 SearXNG 做模糊搜索，输出紧凑 Markdown 表格。
- `read`：直接读取网页真实内容，输出净化后的 `markdown`（默认）或 `text`。

## `--help` 文案要求
### 根命令
- 说明工具支持 `search` 与 `read` 两种流程。
- 给出推荐顺序：先 `search` 再 `read`。
- 若用户已提供 URL：直接调用 `read`。

### `search` 命令
- 说明 `search` 仅返回搜索结果（标题/摘要/链接）。
- 给出参考输出结构：
  - `# | title | url | content | template`

### `read` 命令
- 说明默认输出为 `markdown`，并且会去除杂乱 HTML 标签残留。
- 说明 `--format text` 输出纯文本。
- 给出参考输出示例（markdown 与 text）。

## API 约束（search）
仅使用：

```bash
curl 'https://searx.example.org/search?q=searxng&format=json'
```

实现：
- `GET {server}/search`
- 查询参数仅允许：`q=<query>` 与 `format=json`
- MVP 不使用其他 SearXNG 参数

默认服务器：`https://searx.example.org/`

## 输出约束
### search（stdout）
- 始终输出 Markdown 表格（无 JSON 模式）
- 列固定：`# | title | url | content | template`
- `title/url/content` 来自 `results[]`
- `template` 由 Go 模板渲染（单行）
- 默认模板：`Title={{.Title}} URL={{.URL}} Content={{.Content}}`
- 规范化：
  - 将 `\n`、`\r`、`\t` 替换为空格并折叠重复空格
  - 将单元格内 `|` 转义为 `\|`
  - 推荐截断：title 80 / url 120 / content 160 / template 200

### read（stdout）
- 默认：净化后的 `markdown`
- `--format text`：净化后的纯文本
- 两种模式都必须剔除杂乱 HTML 标签残留

### 错误（stderr）
- 使用结构化 `key=value` 单行格式
- 失败退出码为 `1`

## 命令（Cobra）
- `searxng-cli search <query>`
  - `--limit int`（默认 10）
  - `--template string`
- `searxng-cli read <url>`
  - `--format string`（`markdown|text`，默认 `markdown`）
  - `--timeout duration`
  - `--respect-robots`
  - `--max-bytes int`
  - `--retry int`
- `searxng-cli config init`
- `searxng-cli config view`
- `searxng-cli config use-context <name>`
- `searxng-cli version`

## 认证注入（search/read 共用）
- 支持模式：`none | api_key | basic`
- 支持配置来源：flags / env / config
- 认证输入必须是 base64 编码，CLI 先解码再发送

全局认证 flags：
- `--auth-mode`
- `--auth-header`
- `--auth-api-key`
- `--auth-username`
- `--auth-password`

## 配置
默认路径：`~/.config/searxng-cli/config.yml`

优先级：命令行标志 > 环境变量 > 配置文件 > 默认值

配置示例：

```yaml
apiVersion: searxng-cli/v1
kind: Config
current-context: default
contexts:
  default:
    server: "https://searx.example.org/"
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

## 验收标准
- `searxng-cli search "searxng"` 调用：`GET {server}/search?q=searxng&format=json`
- `search` 的 stdout 为有效 Markdown 表格，且仅包含指定列
- `read` 默认输出净化后的 markdown，`--format text` 输出净化后的纯文本
- 根命令帮助包含流程说明：先 `search` 再 `read`；已知 URL 直接 `read`
