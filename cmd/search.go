/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strings"

	"searxng-cli/internal/apperr"
	"searxng-cli/internal/auth"
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
	Use:   "search <query>",
	Short: "Run fuzzy search and print a Markdown table",
	Long: `This command is only for fuzzy search results (title/summary/link).
To get real page content, use agent-browser, playwright mcp, or any browser automation tool.
Errors are printed to stderr as key=value fields with stable code and retryable flags.
Authentication inputs must be base64-encoded; CLI decodes them before sending requests.`,
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
			return &apperr.ConfigError{Err: err}
		}

		query := strings.Join(args, " ")
		results, err := search.Fetch(cmd.Context(), eff.Server, query, eff.Timeout, auth.Options{
			Mode:         eff.AuthMode,
			APIKeyHeader: eff.AuthHeader,
			APIKey:       eff.AuthAPIKey,
			Username:     eff.AuthUser,
			Password:     eff.AuthPass,
		})
		if err != nil {
			return apperr.Annotate(err, map[string]string{
				"server":  eff.Server,
				"query":   query,
				"timeout": eff.Timeout.String(),
			})
		}

		table, err := render.MarkdownTable(results, eff.Template, eff.Limit)
		if err != nil {
			return &apperr.TemplateError{Err: err}
		}

		_, err = fmt.Fprint(cmd.OutOrStdout(), table)
		return err
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.Flags().IntVar(&searchLimit, "limit", config.DefaultLimit, "Local output limit")
	searchCmd.Flags().StringVar(&searchTemplate, "template", config.DefaultTemplate, "Override row template")
}
