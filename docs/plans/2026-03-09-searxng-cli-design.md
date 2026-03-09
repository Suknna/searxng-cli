# searxng-cli MVP 设计文档

## 目标

实现一个 Go + Cobra 的 CLI，仅通过 SearXNG 返回模糊搜索结果，并将结果输出为紧凑 Markdown 表格。工具不抓取网页正文，不执行 JS。

## 约束

- 仅调用 `GET {server}/search?q=<query>&format=json`。
- 不使用其他 SearXNG 参数。
- `stdout` 永远输出 Markdown 表格，列固定为：`# | title | url | content | template`。
- 错误输出到 `stderr`，请求非 2xx 或解析失败时返回退出码 1。

## 命令与配置

- 命令：`search`、`config init`、`config view`、`config use-context`、`version`。
- 全局标志：`--config`、`--context`、`--server`、`--timeout`、`--verbose`。
- `search` 标志：`--limit`、`--template`。
- 默认配置路径：`~/.config/searxng-cli/config.yml`。
- 配置优先级：命令行标志 > 配置文件 > 默认值。

## 模块设计

- `cmd/`：命令路由与参数解析。
- `internal/config`：配置结构、默认配置、合并优先级、配置文件读写。
- `internal/search`：SearXNG 请求与响应解析。
- `internal/render`：模板渲染、单元格规范化、Markdown 表格输出。

## 输出规则

- template 默认值：`Title={{.Title}} URL={{.URL}} Content={{.Content}}`。
- 表格单元格规范化：
  - 将 `\n`、`\r`、`\t` 替换为空格。
  - 折叠连续空白。
  - `|` 转义为 `\|`。
- 截断策略（推荐值）：title 80、url 120、content 160、template 200。

## 帮助文案要求

根命令与 `search --help` 都必须明确：

- 仅用于模糊搜索（标题/摘要/链接）。
- 不下载/解析/渲染网页，不执行 JS。
- 获取真实页面内容请使用 `agent-browser`、`playwright mcp` 或其他浏览器工具。

## 验证策略

- `internal/search`：`httptest` 断言路径与参数仅有 `q` 和 `format=json`。
- `internal/render`：断言列顺序、换行折叠、竖线转义、模板渲染与截断。
- `internal/config`：断言 flags > file > default 的优先级。
- `cmd`：断言 root/search help 文案包含范围限制。
