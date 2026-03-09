package cmd

import (
	"fmt"
	"strings"
	"time"

	"searxng-cli/internal/read"

	"github.com/spf13/cobra"
)

var (
	readFormat        string
	readTimeout       time.Duration
	readRespectRobots bool
	readMaxBytes      int64
	readRetry         int
)

var readCmd = &cobra.Command{
	Use:   "read <url>",
	Short: "Fetch webpage content as clean markdown or text",
	Long: `Read fetches real webpage content without browser automation.
Default format is markdown, and markdown output is sanitized to remove noisy HTML tags.
Authentication inputs from global auth flags must be base64-encoded.
Reference markdown output:
# Title
Paragraph one.

Reference text output:
Title Paragraph one.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		out, err := read.ReadURL(cmd.Context(), args[0], read.Options{
			Format:        strings.ToLower(readFormat),
			Timeout:       readTimeout,
			RespectRobots: readRespectRobots,
			MaxBytes:      readMaxBytes,
			Retry:         readRetry,
			UserAgent:     "searxng-cli/1.0",
		})
		if err != nil {
			return err
		}
		_, err = fmt.Fprint(cmd.OutOrStdout(), out)
		return err
	},
}

func init() {
	rootCmd.AddCommand(readCmd)
	readCmd.Flags().StringVar(&readFormat, "format", "markdown", "Output format: markdown|text")
	readCmd.Flags().DurationVar(&readTimeout, "timeout", 10*time.Second, "Request timeout")
	readCmd.Flags().BoolVar(&readRespectRobots, "respect-robots", true, "Respect robots.txt rules")
	readCmd.Flags().Int64Var(&readMaxBytes, "max-bytes", 2*1024*1024, "Maximum bytes to read from response")
	readCmd.Flags().IntVar(&readRetry, "retry", 1, "Retry count on fetch failures")
}
