/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"searxng-cli/internal/config"

	"github.com/spf13/cobra"
)

// useContextCmd represents the useContext command
var useContextCmd = &cobra.Command{
	Use:   "use-context <name>",
	Short: "Set current context",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := cfgFile
		if path == "" {
			var err error
			path, err = config.DefaultConfigPath()
			if err != nil {
				return err
			}
		}
		cfg, err := config.LoadRaw(path)
		if err != nil {
			return err
		}
		if _, ok := cfg.Contexts[args[0]]; !ok {
			return fmt.Errorf("context %q not found", args[0])
		}
		cfg.CurrentContext = args[0]
		if err := config.Save(path, cfg); err != nil {
			return err
		}
		_, err = fmt.Fprintln(cmd.OutOrStdout(), args[0])
		return err
	},
}

func init() {
	configCmd.AddCommand(useContextCmd)
}
