/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"time"

	"searxng-cli/internal/config"

	"github.com/spf13/cobra"
)

var (
	cfgFile      string
	flagContext  string
	flagServer   string
	flagTimeout  time.Duration
	flagVerbose  bool
	buildVersion = "dev"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "searxng-cli",
	Short: "通过 SearXNG 获取模糊搜索结果并输出 Markdown 表格",
	Long: `此工具仅用于获取模糊搜索结果（标题/摘要/链接）。
它不下载/解析/渲染网页，也不执行 JS。
要获取真实页面内容：请使用 agent-browser、playwright mcp 或其他任意浏览器工具。`,
	SilenceUsage:  true,
	SilenceErrors: true,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		_, _ = os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "配置文件路径")
	rootCmd.PersistentFlags().StringVar(&flagContext, "context", "", "覆盖当前上下文")
	rootCmd.PersistentFlags().StringVar(&flagServer, "server", "", "覆盖 server")
	rootCmd.PersistentFlags().DurationVar(&flagTimeout, "timeout", 0, "覆盖超时时间，例如 10s")
	rootCmd.PersistentFlags().BoolVar(&flagVerbose, "verbose", false, "输出调试信息")
}

func globalOverrides() config.Overrides {
	return config.Overrides{
		ConfigPath: cfgFile,
		Context:    flagContext,
		Server:     flagServer,
		Timeout:    flagTimeout,
	}
}
