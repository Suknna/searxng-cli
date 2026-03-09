# searxng-cli MVP 规格说明（面向 LLM，仅限 Markdown 表格）
## 目的
一个用于 LLM/智能体的 Go（Cobra）CLI：仅通过 SearXNG 进行模糊搜索，输出紧凑的 Markdown 表格。它**不**获取真实的网页内容。
## 必须在 `--help` 中显示（根命令 + search 命令）
- 此工具仅用于获取模糊搜索结果（标题/摘要/链接）。
- 它不下载/解析/渲染网页，也不执行 JS。
- 要获取真实页面内容：请使用 `agent-browser` 、 `playwright mcp`或其他任意浏览器工具。
## API（硬性约束）
仅使用：
```bash
curl 'https://searxng.searxng.orb.local/search?q=searxng&format=json'
实现：

GET {server}/search
查询参数：q=<query> 和 format=json
MVP 中不使用其他 SearXNG 参数/标志。
测试服务器默认值：https://searxng.searxng.orb.local

输出（仅 stdout，仅 Markdown 表格）
始终向 stdout 输出一个 Markdown 表格。没有 JSON 输出模式。
列（精确）：# | title | url | content | template
对于每个结果行：
title、url、content 来自 SearXNG JSON 的 results[]
template 由 Go 模板渲染（单行；换行符折叠为空格）
默认 template（单行）：

<GOTEMPLATE>
Title={{.Title}} URL={{.URL}} Content={{.Content}}
规范化规则：

将所有单元格中的 \n、\r、\t 替换为空格；折叠重复的空格。
将单元格内的 | 转义为 \|（Markdown 安全）。
可选截断（推荐）：标题 80 字符，url 120 字符，内容 160 字符，模板 200 字符。
错误输出到 stderr；非 2xx 状态码或解析失败时退出码为 1。

命令（Cobra）
searxng-cli search <查询词>
标志：
--limit int（默认 10；本地输出上限）
--template string（覆盖 Go 模板）
searxng-cli config init（将默认配置写入用户配置目录）
searxng-cli config view（打印生效的配置）
searxng-cli config use-context <名称>（设置当前上下文）
searxng-cli version
配置注入（kubectl 风格）
默认配置路径：~/.config/searxng-cli/config.yml（初始化时创建目录）

优先级：命令行标志 > 配置文件 > 默认值

架构：

<YAML>
apiVersion: searxng-cli/v1
kind: Config
current-context: default
contexts:
  default:
    server: "https://searxng.searxng.orb.local"
    timeout: "10s"
    limit: 10
    template: "Title={{.Title}} URL={{.URL}} Content={{.Content}}"
全局标志：

--config string
--context string
--server string
--timeout duration
--verbose
验收标准
searxng-cli search "searxng" 调用：GET {server}/search?q=searxng&format=json
stdout 是一个有效的 Markdown 表格，且仅包含指定的列。
帮助文本包含“仅用于模糊搜索；获取页面内容请使用 agent-browser/playwright mcp或者其他浏览器工具”的提示。