/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strings"

	"searxng-cli/internal/config"
	"searxng-cli/internal/render"
	"searxng-cli/internal/search"

	"github.com/spf13/cobra"
)

var (
	searchLimit    int
	searchTemplate string
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search <查询词>",
	Short: "执行模糊搜索并输出 Markdown 表格",
	Long: `此命令仅用于模糊搜索结果（标题/摘要/链接）。
要获取真实页面内容：请使用 agent-browser、playwright mcp 或其他任意浏览器工具。`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		o := globalOverrides()
		if cmd.Flags().Changed("limit") {
			o.Limit = &searchLimit
		}
		if cmd.Flags().Changed("template") {
			o.Template = &searchTemplate
		}

		eff, _, err := config.LoadEffective(o)
		if err != nil {
			return err
		}
		if flagVerbose {
			_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "server=%s context=%s timeout=%s\n", eff.Server, eff.ContextName, eff.Timeout)
		}

		query := strings.Join(args, " ")
		results, err := search.Fetch(cmd.Context(), eff.Server, query, eff.Timeout)
		if err != nil {
			return err
		}

		table, err := render.MarkdownTable(results, eff.Template, eff.Limit)
		if err != nil {
			return err
		}

		_, err = fmt.Fprint(cmd.OutOrStdout(), table)
		return err
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.Flags().IntVar(&searchLimit, "limit", config.DefaultLimit, "本地输出上限")
	searchCmd.Flags().StringVar(&searchTemplate, "template", config.DefaultTemplate, "覆盖模板")
}
