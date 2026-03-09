/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"time"

	"searxng-cli/internal/apperr"
	"searxng-cli/internal/config"

	"github.com/spf13/cobra"
)

var (
	cfgFile      string
	flagContext  string
	flagServer   string
	flagTimeout  time.Duration
	flagAuthMode *string
	flagAuthHdr  *string
	flagAuthKey  *string
	flagAuthUser *string
	flagAuthPass *string
	buildVersion = "dev"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "searxng-cli",
	Short: "Get fuzzy search results from SearXNG as a Markdown table",
	Long: `This tool is only for fuzzy search results (title/summary/link).
It does not download, parse, or render web pages, and it does not execute JS.
Recommended workflow: search first, then read.
If you already have a URL, call read directly.`,
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
		line := renderCLIError(err)
		_, _ = os.Stderr.WriteString(line + "\n")
		os.Exit(1)
	}
}

func renderCLIError(err error) string {
	return apperr.RenderKV(apperr.FromError(err))
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Path to config file")
	rootCmd.PersistentFlags().StringVar(&flagContext, "context", "", "Override active context")
	rootCmd.PersistentFlags().StringVar(&flagServer, "server", "", "Override server")
	rootCmd.PersistentFlags().DurationVar(&flagTimeout, "timeout", 0, "Override timeout (for example 10s)")
	flagAuthMode = rootCmd.PersistentFlags().String("auth-mode", "", "Auth mode: none|api_key|basic")
	flagAuthHdr = rootCmd.PersistentFlags().String("auth-header", "", "Auth header name for api_key mode")
	flagAuthKey = rootCmd.PersistentFlags().String("auth-api-key", "", "Base64 API key value")
	flagAuthUser = rootCmd.PersistentFlags().String("auth-username", "", "Base64 username for basic mode")
	flagAuthPass = rootCmd.PersistentFlags().String("auth-password", "", "Base64 password for basic mode")
}

func globalOverrides() config.Overrides {
	return config.Overrides{
		ConfigPath: cfgFile,
		Context:    flagContext,
		Server:     flagServer,
		Timeout:    flagTimeout,
		AuthMode:   strPtrIfNotEmpty(flagAuthMode),
		AuthHeader: strPtrIfNotEmpty(flagAuthHdr),
		AuthAPIKey: strPtrIfNotEmpty(flagAuthKey),
		AuthUser:   strPtrIfNotEmpty(flagAuthUser),
		AuthPass:   strPtrIfNotEmpty(flagAuthPass),
	}
}

func strPtrIfNotEmpty(v *string) *string {
	if v == nil || *v == "" {
		return nil
	}
	return v
}
